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

	Post interface {
		CreatePost(ctx context.Context, in *PostReq, opts ...grpc.CallOption) (*EmptyResp, error)
		UpdatePost(ctx context.Context, in *PostReq, opts ...grpc.CallOption) (*EmptyResp, error)
		DeletePost(ctx context.Context, in *DeletesReq, opts ...grpc.CallOption) (*EmptyResp, error)
		PagePost(ctx context.Context, in *PostQueryReq, opts ...grpc.CallOption) (*PostListResp, error)
	}

	defaultPost struct {
		cli zrpc.Client
	}
)

func NewPost(cli zrpc.Client) Post {
	return &defaultPost{
		cli: cli,
	}
}

func (m *defaultPost) CreatePost(ctx context.Context, in *PostReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	client := pb.NewPostClient(m.cli.Conn())
	return client.CreatePost(ctx, in, opts...)
}

func (m *defaultPost) UpdatePost(ctx context.Context, in *PostReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	client := pb.NewPostClient(m.cli.Conn())
	return client.UpdatePost(ctx, in, opts...)
}

func (m *defaultPost) DeletePost(ctx context.Context, in *DeletesReq, opts ...grpc.CallOption) (*EmptyResp, error) {
	client := pb.NewPostClient(m.cli.Conn())
	return client.DeletePost(ctx, in, opts...)
}

func (m *defaultPost) PagePost(ctx context.Context, in *PostQueryReq, opts ...grpc.CallOption) (*PostListResp, error) {
	client := pb.NewPostClient(m.cli.Conn())
	return client.PagePost(ctx, in, opts...)
}