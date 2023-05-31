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

type UpdateSysApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSysApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSysApplicationLogic {
	return &UpdateSysApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSysApplicationLogic) UpdateSysApplication(in *pb.UpdateSysApplicationReq) (*pb.UpdateSysApplicationResp, error) {
	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, in.Id)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	_ = copier.Copiers(app, in)
	err := l.svcCtx.AppModel.Update(l.ctx, app)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db update app error: %+v", err)
	}

	return &pb.UpdateSysApplicationResp{}, nil
}
