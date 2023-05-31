// Code generated by goctl. DO NOT EDIT.
// Source: koala.proto

package server

import (
	"context"

	"koala/gateway/internal/logic/api"
	"koala/gateway/internal/svc"
	"koala/gateway/pb"
)

type ApiServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedApiServer
}

func NewApiServer(svcCtx *svc.ServiceContext) *ApiServer {
	return &ApiServer{
		svcCtx: svcCtx,
	}
}

func (s *ApiServer) CreateApi(ctx context.Context, in *pb.ApiReq) (*pb.EmptyResp, error) {
	l := apilogic.NewCreateApiLogic(ctx, s.svcCtx)
	return l.CreateApi(in)
}

func (s *ApiServer) UpdateApi(ctx context.Context, in *pb.ApiReq) (*pb.EmptyResp, error) {
	l := apilogic.NewUpdateApiLogic(ctx, s.svcCtx)
	return l.UpdateApi(in)
}

func (s *ApiServer) DeleteApi(ctx context.Context, in *pb.DeletesReq) (*pb.EmptyResp, error) {
	l := apilogic.NewDeleteApiLogic(ctx, s.svcCtx)
	return l.DeleteApi(in)
}

func (s *ApiServer) PageApi(ctx context.Context, in *pb.ApiQueryReq) (*pb.ApiListResp, error) {
	l := apilogic.NewPageApiLogic(ctx, s.svcCtx)
	return l.PageApi(in)
}