// Code generated by goctl. DO NOT EDIT.
// Source: koala.proto

package client

import (
	"context"

	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AccountListResp          = pb.AccountListResp
	AccountQueryReq          = pb.AccountQueryReq
	AccountReq               = pb.AccountReq
	ApiListResp              = pb.ApiListResp
	ApiQueryReq              = pb.ApiQueryReq
	ApiReq                   = pb.ApiReq
	BoolRep                  = pb.BoolRep
	CreateSysApplicationReq  = pb.CreateSysApplicationReq
	CreateSysApplicationResp = pb.CreateSysApplicationResp
	DeleteSysApplicationReq  = pb.DeleteSysApplicationReq
	DeleteSysApplicationResp = pb.DeleteSysApplicationResp
	DeletesReq               = pb.DeletesReq
	DeptListResp             = pb.DeptListResp
	DeptQueryReq             = pb.DeptQueryReq
	DeptReq                  = pb.DeptReq
	DetailSysApplicationReq  = pb.DetailSysApplicationReq
	DetailSysApplicationResp = pb.DetailSysApplicationResp
	EmptyReq                 = pb.EmptyReq
	EmptyResp                = pb.EmptyResp
	GenerateTokenReq         = pb.GenerateTokenReq
	GenerateTokenResp        = pb.GenerateTokenResp
	GetUserInfoReq           = pb.GetUserInfoReq
	GetUserInfoResp          = pb.GetUserInfoResp
	IdReq                    = pb.IdReq
	IdsResp                  = pb.IdsResp
	LoginReq                 = pb.LoginReq
	LoginResp                = pb.LoginResp
	MenuQueryReq             = pb.MenuQueryReq
	MenuReq                  = pb.MenuReq
	MenusListResp            = pb.MenusListResp
	MenusListResp_MenuMeta   = pb.MenusListResp_MenuMeta
	MenusResp                = pb.MenusResp
	Permission               = pb.Permission
	PermsReq                 = pb.PermsReq
	PermsResp                = pb.PermsResp
	PolicyReq                = pb.PolicyReq
	PostListResp             = pb.PostListResp
	PostQueryReq             = pb.PostQueryReq
	PostReq                  = pb.PostReq
	QueryReq                 = pb.QueryReq
	RegisterReq              = pb.RegisterReq
	RegisterResp             = pb.RegisterResp
	ResetPasswordReq         = pb.ResetPasswordReq
	RoleBindReq              = pb.RoleBindReq
	RoleListResp             = pb.RoleListResp
	RolePermission           = pb.RolePermission
	RoleQueryReq             = pb.RoleQueryReq
	RoleReq                  = pb.RoleReq
	RuleResp                 = pb.RuleResp
	RulesReq                 = pb.RulesReq
	SearchSysApplicationReq  = pb.SearchSysApplicationReq
	SearchSysApplicationResp = pb.SearchSysApplicationResp
	SysApplicationView       = pb.SysApplicationView
	UpdatePasswordReq        = pb.UpdatePasswordReq
	UpdateSysApplicationReq  = pb.UpdateSysApplicationReq
	UpdateSysApplicationResp = pb.UpdateSysApplicationResp
	User                     = pb.User
	UserMenusResp            = pb.UserMenusResp
	UserPermsVerifyReq       = pb.UserPermsVerifyReq
	UserReq                  = pb.UserReq

	Api interface {
		CreateApi(ctx context.Context, in *ApiReq, opts ...grpc.CallOption) (*EmptyResp, error)
		UpdateApi(ctx context.Context, in *ApiReq, opts ...grpc.CallOption) (*EmptyResp, error)
		DeleteApi(ctx context.Context, in *DeletesReq, opts ...grpc.CallOption) (*EmptyResp, error)
		PageApi(ctx context.Context, in *ApiQueryReq, opts ...grpc.CallOption) (*ApiListResp, error)
	}

	defaultApi struct {
		cli zrpc.Client
	}
)

func NewApi(cli zrpc.Client) Api {
	return &defaultApi{
		cli: cli,
	}
}

func (m *defaultApi) CreateApi(ctx context.Context, in *ApiReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	client := pb.NewApiClient(m.cli.Conn())
	return client.CreateApi(ctx, in, opts...)
}

func (m *defaultApi) UpdateApi(ctx context.Context, in *ApiReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	client := pb.NewApiClient(m.cli.Conn())
	return client.UpdateApi(ctx, in, opts...)
}

func (m *defaultApi) DeleteApi(ctx context.Context, in *DeletesReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	client := pb.NewApiClient(m.cli.Conn())
	return client.DeleteApi(ctx, in, opts...)
}

func (m *defaultApi) PageApi(ctx context.Context, in *ApiQueryReq, opts ...grpc.CallOption) (*ApiListResp, error) {
	client := pb.NewApiClient(m.cli.Conn())
	return client.PageApi(ctx, in, opts...)
}
