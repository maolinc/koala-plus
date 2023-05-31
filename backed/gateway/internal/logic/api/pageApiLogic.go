package apilogic

import (
	"context"
	"encoding/json"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageApiLogic {
	return &PageApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageApiLogic) PageApi(in *pb.ApiQueryReq) (*pb.ApiListResp, error) {
	cond := model.SysApiQuery{}
	_ = copier.Copiers(&cond, &in)
	_ = json.Unmarshal([]byte(in.SearchPlus), &cond.SearchPlus)

	total, list, err := l.svcCtx.SysApiModel.FindByPage(l.ctx, &cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	resp := pb.ApiListResp{Total: total}
	_ = copier.Copiers(&resp.List, list)

	return &resp, nil
}
