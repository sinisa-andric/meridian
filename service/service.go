package service

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/c12s/meridian/helper"
	"github.com/c12s/meridian/model"
	"github.com/c12s/meridian/storage"
	cPb "github.com/c12s/scheme/celestial"
	mPb "github.com/c12s/scheme/meridian"
	sg "github.com/c12s/stellar-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	instrument map[string]string
	db         storage.DB
	apollo     string
	//meridian   string
}

func (s *Server) List(ctx context.Context, req *cPb.ListReq) (*cPb.ListResp, error) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.List")
	defer span.Finish()
	fmt.Println(span)

	token, err := helper.ExtractToken(ctx)
	if err != nil {
		span.AddLog(&sg.KV{Key: "token error", Value: err.Error()})
		return nil, err
	}

	err = s.auth(ctx, listOpt(req, token))
	if err != nil {
		span.AddLog(&sg.KV{Key: "auth error", Value: err.Error()})
		return nil, err
	}

	_, err = s.checkNS(ctx, req.Extras["user"], req.Extras["namespace"])
	if err != nil {
		span.AddLog(&sg.KV{Key: "check ns error", Value: err.Error()})
		return nil, err
	}

	err, rsp := s.db.List(
		helper.AppendToken(
			sg.NewTracedGRPCContext(ctx, span),
			token,
		),
		req.Extras,
	)
	if err != nil {
		span.AddLog(&sg.KV{Key: "roles list error", Value: err.Error()})
		return nil, err
	}
	return rsp, nil
}

func (s *Server) Mutate(ctx context.Context, req *cPb.MutateReq) (*cPb.MutateResp, error) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.Mutate")
	defer span.Finish()
	fmt.Println(span)

	token, err := helper.ExtractToken(ctx)
	if err != nil {
		span.AddLog(&sg.KV{Key: "token error", Value: err.Error()})
		return nil, err
	}
	fmt.Println("{{TOKEN MUTATE MERIDIAN}}", token)

	err = s.auth(ctx, mutateOpt(req, token))
	if err != nil {
		span.AddLog(&sg.KV{Key: "auth error", Value: err.Error()})
		return nil, err
	}

	ns, err := s.checkNS(ctx, req.Mutate.UserId, req.Mutate.Extras["namespace"])
	if err != nil {
		span.AddLog(&sg.KV{Key: "namespace check error", Value: err.Error()})
		return nil, err
	}
	if ns != "" {
		span.AddLog(&sg.KV{Key: "check ns error", Value: "Namespace already exis"})
		return nil, errors.New("namespace already exis")
	}

	err, rsp := s.db.Mutate(helper.AppendToken(
		sg.NewTracedGRPCContext(ctx, span),
		token,
	), req)
	if err != nil {
		span.AddLog(&sg.KV{Key: "namespace mutate error", Value: err.Error()})
		return nil, err
	}
	return rsp, nil
}

func (s *Server) Exists(ctx context.Context, req *mPb.NSReq) (*mPb.NSResp, error) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.Exists")
	defer span.Finish()
	fmt.Println(span)

	err, rsp := s.db.Exists(
		sg.NewTracedGRPCContext(ctx, span), req,
	)
	if err != nil {
		fmt.Println("{{service.Exists err != nil}}")
		span.AddLog(&sg.KV{Key: "namespace exists error", Value: err.Error()})
		return nil, err
	}

	fmt.Println("{{service.Exists err == nil}}")
	return rsp, nil
}

func (s *Server) Delete(ctx context.Context, req *mPb.NSReq) (*mPb.NSResp, error) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.Delete")
	defer span.Finish()
	fmt.Println(span)

	token, err := helper.ExtractToken(ctx)
	if err != nil {
		span.AddLog(&sg.KV{Key: "token error", Value: err.Error()})
		return nil, err
	}

	// check this later
	// _, err = s.checkNS(ctx, req.Mutate.UserId, req.Mutate.Namespace)
	// if err != nil {
	// 	span.AddLog(&sg.KV{Key: "check ns error", Value: err.Error()})
	// 	return nil, err
	// }

	err, rsp := s.db.Delete(
		helper.AppendToken(
			sg.NewTracedGRPCContext(ctx, span),
			token,
		), req,
	)

	if err != nil {
		span.AddLog(&sg.KV{Key: "namespace exists error", Value: err.Error()})
		return nil, err
	}

	return rsp, nil
}

func Run(db storage.DB, conf *model.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", conf.Address)
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	meridianServer := &Server{
		instrument: conf.InstrumentConf,
		db:         db,
		apollo:     conf.Apollo,
	}

	n, err := sg.NewCollector(meridianServer.instrument["address"], meridianServer.instrument["stopic"])
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := sg.InitCollector(meridianServer.instrument["location"], n)
	if err != nil {
		fmt.Println(err)
		return
	}
	go c.Start(ctx, 15*time.Second)

	fmt.Println("MeridianService RPC Started")
	mPb.RegisterMeridianServiceServer(server, meridianServer)
	server.Serve(lis)
}
