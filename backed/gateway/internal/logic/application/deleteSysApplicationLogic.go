package applicationlogic

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"koala/gateway/internal/tools/errorx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSysApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSysApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSysApplicationLogic {
	return &DeleteSysApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSysApplicationLogic) DeleteSysApplication(in *pb.DeleteSysApplicationReq) (*pb.DeleteSysApplicationResp, error) {
	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, in.Id)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	err := l.svcCtx.AppModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
		err = l.svcCtx.AppModel.Delete(l.ctx, in.Id)
		if err != nil {
			return errors.Wrapf(errorx.DbError, "db app delete error: %+v", err)
		}
		// clear app all permission and role
		res := l.svcCtx.CBS.ClearApp(app.Perms)
		if !res {
			return errorx.NewMsg("App has been deleted, and some permissions have not been cleaned up")
		}
		return nil
	})

	return &pb.DeleteSysApplicationResp{}, err
}
