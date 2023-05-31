package postlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PagePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPagePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PagePostLogic {
	return &PagePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PagePostLogic) PagePost(in *pb.PostQueryReq) (*pb.PostListResp, error) {
	cond := &model.SysPostQuery{}
	_ = copier.Copiers(&cond, &in)

	total, list, err := l.svcCtx.SysPostModel.FindByPage(l.ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	resp := pb.PostListResp{Total: total}
	_ = copier.Copiers(&resp.List, list)

	return &resp, nil
}
