package applicationlogic

import (
	"context"
	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"koala/gateway/internal/tools/errorx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailSysApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailSysApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailSysApplicationLogic {
	return &DetailSysApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailSysApplicationLogic) DetailSysApplication(in *pb.DetailSysApplicationReq) (*pb.DetailSysApplicationResp, error) {
	data, err := l.svcCtx.AppModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db app select error: %+v", err)
	}

	resp := &pb.DetailSysApplicationResp{}
	_ = copier.Copiers(resp, data)

	return resp, nil
}
