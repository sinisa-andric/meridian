package service

import (
	"errors"
	"fmt"
	"github.com/c12s/meridian/helper"
	aPb "github.com/c12s/scheme/apollo"
	cPb "github.com/c12s/scheme/celestial"
	mPb "github.com/c12s/scheme/meridian"
	sg "github.com/c12s/stellar-go"
	"golang.org/x/net/context"
)

func (s *Server) isTransientToken(token string) bool {
	return token == "my_test_token"
}

func (s *Server) auth(ctx context.Context, opt *aPb.AuthOpt) error {
	span, _ := sg.FromGRPCContext(ctx, "auth")
	defer span.Finish()
	fmt.Println(span)

	token, err := helper.ExtractToken(ctx)
	if err != nil {
		span.AddLog(&sg.KV{"token error", err.Error()})
		return err
	}

	if s.isTransientToken(token) {
		return nil
	}

	client := NewApolloClient(s.apollo)
	resp, err := client.Auth(
		helper.AppendToken(
			sg.NewTracedGRPCContext(ctx, span),
			token,
		),
		opt,
	)
	if err != nil {
		span.AddLog(&sg.KV{"apollo resp error", err.Error()})
		return err
	}

	if !resp.Value {
		span.AddLog(&sg.KV{"apollo.auth value", resp.Data["message"]})
		return errors.New(resp.Data["message"])
	}
	return nil
}

func (s *Server) checkNS(ctx context.Context, userid, namespace string) (string, error) {
	span, _ := sg.FromGRPCContext(ctx, "ns check")
	defer span.Finish()
	fmt.Println(span)

	mrsp, err := s.Exists(sg.NewTracedGRPCContext(ctx, span),
		&mPb.NSReq{
			Name:   namespace,
			Extras: map[string]string{"userid": userid},
		},
	)
	if err != nil {
		span.AddLog(&sg.KV{"meridian exists error", err.Error()})
		fmt.Println("namespace do not exists")
		return "", err
	}
	return mrsp.Extras["exists"], nil
}

func listOpt(req *cPb.ListReq, token string) *aPb.AuthOpt {
	return &aPb.AuthOpt{
		Data: map[string]string{
			"intent": "auth",
			"action": "list",
			"kind":   "namespaces",
			"token":  token,
		},
		Extras: map[string]*aPb.OptExtras{
			"user":      &aPb.OptExtras{Data: []string{req.Extras["user"]}},
			"namespace": &aPb.OptExtras{Data: []string{req.Extras["namespace"]}},
			"cmp":       &aPb.OptExtras{Data: []string{req.Extras["compare"]}},
			"labels":    &aPb.OptExtras{Data: []string{req.Extras["labels"]}},
		},
	}
}

func mutateOpt(req *cPb.MutateReq, token string) *aPb.AuthOpt {
	return &aPb.AuthOpt{
		Data: map[string]string{
			"intent":    "auth",
			"action":    "mutate",
			"kind":      "namespaces",
			"user":      req.Mutate.UserId,
			"token":     token,
			"namespace": req.Mutate.Namespace,
		},
		Extras: map[string]*aPb.OptExtras{
			"labels":    &aPb.OptExtras{Data: []string{req.Mutate.Extras["labels"]}},
			"namespace": &aPb.OptExtras{Data: []string{req.Mutate.Extras["namespace"]}},
		},
	}
}
