package accountlogic

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

type PageAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageAccountLogic {
	return &PageAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageAccountLogic) PageAccount(in *pb.AccountQueryReq) (*pb.AccountListResp, error) {
	cond := model.SysUserQuery{}
	_ = copier.Copiers(cond, &in)

	total, list, err := l.svcCtx.UserModel.FindByPage(l.ctx, &cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	resp := pb.AccountListResp{Total: total}
	_ = copier.Copiers(&resp.List, list)

	return &resp, nil
}
