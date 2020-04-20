package storage

import (
	"context"
	cPb "github.com/c12s/scheme/celestial"
	mPb "github.com/c12s/scheme/meridian"
)

type DB interface {
	List(ctx context.Context, extras map[string]string) (error, *cPb.ListResp)
	Mutate(ctx context.Context, req *cPb.MutateReq) (error, *cPb.MutateResp)
	Exists(ctx context.Context, req *mPb.NSReq) (error, *mPb.NSResp)
	Delete(ctx context.Context, req *mPb.NSReq) (error, *mPb.NSResp)
	Init() error
}
