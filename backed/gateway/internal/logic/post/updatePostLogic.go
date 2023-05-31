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

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(in *pb.PostReq) (*pb.EmptyResp, error) {
	old, _ := l.svcCtx.SysPostModel.FindOne(l.ctx, in.Id)
	if old == nil || old.Id == 0 {
		return nil, errorx.NewMsg("岗位不存在")
	}

	var post model.SysPost
	_ = copier.Copiers(&post, &in)
	err := l.svcCtx.SysPostModel.Update(l.ctx, &post)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("岗位更新失败"), "db update err:%v", err)
	}

	return &pb.EmptyResp{}, nil
}
