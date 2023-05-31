package postlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *pb.PostReq) (*pb.EmptyResp, error) {
	var post model.SysPost
	_ = copier.Copiers(&post, in)
	post.Id = 0

	err := l.svcCtx.SysPostModel.Insert(l.ctx, &post)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db insert err:%v", err)
	}

	return &pb.EmptyResp{}, nil
}
