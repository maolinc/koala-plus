package applicationlogic

import (
	"context"
	"encoding/json"
	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageSysApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageSysApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageSysApplicationLogic {
	return &PageSysApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageSysApplicationLogic) PageSysApplication(in *pb.SearchSysApplicationReq) (*pb.SearchSysApplicationResp, error) {
	cond := model.SysApplicationQuery{}
	_ = copier.Copiers(&cond, &in)
	_ = json.Unmarshal([]byte(in.SearchPlus), &cond.SearchPlus)

	total, list, err := l.svcCtx.AppModel.FindByPage(l.ctx, &cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db page app error: %+v", err)
	}

	resp := &pb.SearchSysApplicationResp{Total: total}
	_ = copier.Copiers(&resp.List, list)

	return resp, nil
}
