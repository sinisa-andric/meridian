package etcd

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/c12s/meridian/model"
	"github.com/c12s/meridian/storage"
	cPb "github.com/c12s/scheme/celestial"
	rPb "github.com/c12s/scheme/core"
	mPb "github.com/c12s/scheme/meridian"
	sg "github.com/c12s/stellar-go"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/protobuf/proto"
)

type ETCD struct {
	kv     clientv3.KV
	client *clientv3.Client
	cache  storage.Cacher
}

func New(conf *model.Config, cache storage.Cacher, timeout time.Duration) (*ETCD, error) {
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: timeout,
		Endpoints:   conf.DB,
	})

	if err != nil {
		return nil, err
	}

	return &ETCD{
		kv:     clientv3.NewKV(cli),
		client: cli,
		cache:  cache,
	}, nil
}

func (db *ETCD) Close() { db.client.Close() }

func (e *ETCD) get(ctx context.Context, key string) (string, int64, string, string) {
	span, _ := sg.FromGRPCContext(ctx, "get")
	defer span.Finish()
	fmt.Println(span)

	chspan := span.Child("etcd.get")
	gresp, err := e.kv.Get(ctx, key)
	if err != nil {
		chspan.AddLog(&sg.KV{Key: "etcd get error", Value: err.Error()})
		return "", 0, "", ""
	}
	go chspan.Finish()

	for _, item := range gresp.Kvs {
		nsTask := &rPb.Task{}
		err = proto.Unmarshal(item.Value, nsTask.ProtoReflect())
		if err != nil {
			span.AddLog(&sg.KV{Key: "unmarshall etcd get error", Value: err.Error()})
			return "", 0, "", ""
		}
		return nsTask.Namespace, nsTask.Timestamp, nsTask.Extras["namespace"], nsTask.Extras["labels"]
	}
	return "", 0, "", ""
}

func (e *ETCD) List(ctx context.Context, extra map[string]string) (error, *cPb.ListResp) {
	span, _ := sg.FromGRPCContext(ctx, "list")
	defer span.Finish()
	fmt.Println(span)

	name := extra["name"]
	cmp := extra["compare"]
	els := split(extra["labels"])
	user := extra["user"]
	sort.Strings(els)

	datas := []*cPb.Data{}
	if name == "" {
		chspan := span.Child("etcd.get searchLabels")
		gresp, err := e.kv.Get(ctx, NSLabels(user), clientv3.WithPrefix(),
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
		if err != nil {
			chspan.AddLog(&sg.KV{Key: "etcd.get error", Value: err.Error()})
			return err, nil
		}
		go chspan.Finish()

		for _, item := range gresp.Kvs {
			key := string(item.Key)
			newKey := join("/", ssplit(key, "/labels/"))
			ls := SplitLabels(string(item.Value))

			data := &cPb.Data{Data: map[string]string{}}
			switch cmp {
			case "all":
				if len(ls) == len(els) && Compare(ls, els, true) {
					ns, timestamp, name, labels := e.get(sg.NewTracedGRPCContext(ctx, span), newKey)
					if ns != "" {
						data.Data["namespace"] = ns
						data.Data["age"] = strconv.FormatInt(timestamp, 10)
						data.Data["name"] = name
						data.Data["labels"] = labels
					}
					datas = append(datas, data)
				}
			case "any":
				if Compare(ls, els, false) {
					ns, timestamp, name, labels := e.get(sg.NewTracedGRPCContext(ctx, span), newKey)
					if ns != "" {
						data.Data["namespace"] = ns
						data.Data["age"] = strconv.FormatInt(timestamp, 10)
						data.Data["name"] = name
						data.Data["labels"] = labels
					}
					datas = append(datas, data)
				}
			}
		}
	} else {
		data := &cPb.Data{Data: map[string]string{}}
		nsKey := newNSKeyspace(user, name)
		ns, timestamp, name, labels := e.get(sg.NewTracedGRPCContext(ctx, span), nsKey)
		if ns != "" {
			data.Data["namespace"] = ns
			data.Data["age"] = strconv.FormatInt(timestamp, 10)
			data.Data["name"] = name
			data.Data["labels"] = labels
		}
		datas = append(datas, data)
	}
	return nil, &cPb.ListResp{Data: datas}
}

func (e *ETCD) Mutate(ctx context.Context, req *cPb.MutateReq) (error, *cPb.MutateResp) {
	span, _ := sg.FromGRPCContext(ctx, "mutate")
	defer span.Finish()
	fmt.Println(span)

	task := req.Mutate
	namespace := task.Extras["namespace"]
	labels := task.Extras["labels"]

	nsKey := newNSKeyspace(task.UserId, namespace)
	nsData, merr := proto.Marshal(task.ProtoReflect())
	if merr != nil {
		span.AddLog(&sg.KV{Key: "etcd.put key error", Value: merr.Error()})
		return merr, nil
	}

	chspan1 := span.Child("etcd.put key")
	_, err := e.kv.Put(ctx, nsKey, string(nsData))
	if err != nil {
		chspan1.AddLog(&sg.KV{Key: "etcd.put key error", Value: err.Error()})
		return err, nil
	}
	chspan1.Finish()

	chspan2 := span.Child("etcd.put labels")
	lKey := newNSLabelsKeyspace(task.UserId, namespace)
	_, err = e.kv.Put(ctx, lKey, labels)
	if err != nil {
		chspan2.AddLog(&sg.KV{Key: "etcd.put labels error", Value: err.Error()})
		return err, nil
	}
	chspan2.Finish()
	return nil, &cPb.MutateResp{Error: "Namespace added."}
}

func (e *ETCD) Exists(ctx context.Context, req *mPb.NSReq) (error, *mPb.NSResp) {
	span, _ := sg.FromGRPCContext(ctx, "exists")
	defer span.Finish()
	fmt.Println(span)

	userid := req.Extras["userid"]
	key := newNSKeyspace(userid, req.Name)
	fmt.Println("LOOKUP", key)

	_, err := e.cache.Get(key)
	if err == nil {
		fmt.Println("EXISTS CACHE HIT")
		return nil, &mPb.NSResp{
			Extras: map[string]string{
				"exists": req.Name,
			},
		}
	}

	chspan := span.Child("etcd.get")
	gresp, err := e.kv.Get(ctx, key)
	if err != nil {
		chspan.AddLog(&sg.KV{Key: "etcd get error", Value: err.Error()})
		return err, nil
	}
	go chspan.Finish()

	if len(gresp.Kvs) > 0 {
		e.cache.Put(key, true, 10*time.Minute)
		fmt.Println("EXISTS CACHEd")

		return nil, &mPb.NSResp{
			Extras: map[string]string{
				"exists": req.Name,
			},
		}
	}
	fmt.Println("NAMESPACE NOT EXISTING")
	return errors.New(fmt.Sprintf("%s do not exists...", req.Name)), nil
}

func (e *ETCD) Delete(ctx context.Context, req *mPb.NSReq) (error, *mPb.NSResp) {
	span, _ := sg.FromGRPCContext(ctx, "delete")
	defer span.Finish()
	fmt.Println(span)

	userid := req.Extras["userid"]

	key := newNSKeyspace(userid, req.Name)
	chspan := span.Child("etcd.delete")
	gresp, err := e.kv.Delete(ctx, key)
	if err != nil {
		chspan.AddLog(&sg.KV{Key: "etcd delete error", Value: err.Error()})
		return err, nil
	}
	go chspan.Finish()

	lkey := newNSLabelsKeyspace(userid, req.Name)
	chspan1 := span.Child("etcd.delete")
	_, err = e.kv.Delete(ctx, lkey)
	if err != nil {
		chspan1.AddLog(&sg.KV{Key: "etcd delete error", Value: err.Error()})
		return err, nil
	}
	go chspan1.Finish()

	if gresp.Deleted > 0 {
		return nil, &mPb.NSResp{
			Extras: map[string]string{
				"deleted": req.Name,
			},
		}
	}

	return nil, &mPb.NSResp{
		Extras: map[string]string{
			"deleted": "",
		},
	}
}
