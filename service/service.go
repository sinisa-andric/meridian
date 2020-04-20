package service

import (
	"errors"
	"fmt"
	"github.com/c12s/meridian/helper"
	"github.com/c12s/meridian/model"
	"github.com/c12s/meridian/storage"
	aPb "github.com/c12s/scheme/apollo"
	cPb "github.com/c12s/scheme/celestial"
	mPb "github.com/c12s/scheme/meridian"
	sg "github.com/c12s/stellar-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Server struct {
	instrument map[string]string
	db         storage.DB
	apollo     string
}

func (s *Server) auth(ctx context.Context) (string, error) {
	span, _ := sg.FromGRPCContext(ctx, "auth")
	defer span.Finish()
	fmt.Println(span)

	token, err := helper.ExtractToken(ctx)
	if err != nil {
		span.AddLog(&sg.KV{"token error", err.Error()})
		return "", err
	}

	client := NewApolloClient(s.apollo)
	resp, err := client.Auth(
		helper.AppendToken(
			sg.NewTracedGRPCContext(ctx, span),
			token,
		),
		&aPb.AuthOpt{
			Data: map[string]string{"intent": "auth"},
		},
	)
	if err != nil {
		span.AddLog(&sg.KV{"apollo resp error", err.Error()})
		return "", err
	}

	if !resp.Value {
		span.AddLog(&sg.KV{"apollo.auth value", resp.Data["message"]})
		return "", errors.New(resp.Data["message"])
	}

	return token, nil
}

func (s *Server) List(ctx context.Context, req *cPb.ListReq) (*cPb.ListResp, error) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.List")
	defer span.Finish()
	fmt.Println(span)

	token, err := s.auth(ctx)
	if err != nil {
		span.AddLog(&sg.KV{"apollo resp error", err.Error()})
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
		span.AddLog(&sg.KV{"roles list error", err.Error()})
		return nil, err
	}
	return rsp, nil
}

func (s *Server) Mutate(ctx context.Context, req *cPb.MutateReq) (*cPb.MutateResp, error) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.Mutate")
	defer span.Finish()
	fmt.Println(span)

	token, err := s.auth(ctx)
	if err != nil {
		span.AddLog(&sg.KV{"apollo resp error", err.Error()})
		return nil, err
	}

	err, rsp := s.db.Mutate(helper.AppendToken(
		sg.NewTracedGRPCContext(ctx, span),
		token,
	), req,
	)
	if err != nil {
		span.AddLog(&sg.KV{"namespace mutate error", err.Error()})
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
		span.AddLog(&sg.KV{"namespace exists error", err.Error()})
		return nil, err
	}
	return rsp, err
}

func (s *Server) Delete(ctx context.Context, req *mPb.NSReq) (error, *mPb.NSResp) {
	span, _ := sg.FromGRPCContext(ctx, "meridian.Delete")
	defer span.Finish()
	fmt.Println(span)

	token, err := s.auth(ctx)
	if err != nil {
		span.AddLog(&sg.KV{"apollo resp error", err.Error()})
		return err, nil
	}

	err, rsp := s.db.Delete(
		helper.AppendToken(
			sg.NewTracedGRPCContext(ctx, span),
			token,
		), req,
	)
	if err != nil {
		span.AddLog(&sg.KV{"namespace exists error", err.Error()})
		return err, nil
	}
	return nil, rsp
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

	err = db.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("MeridianService RPC Started")
	mPb.RegisterMeridianServiceServer(server, meridianServer)
	server.Serve(lis)
}
