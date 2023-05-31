package apilogic

import (
	"context"
	"gorm.io/gorm"
	"koala/gateway/internal/tools/errorx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteApiLogic) DeleteApi(in *pb.DeletesReq) (*pb.EmptyResp, error) {
	for _, id := range in.Ids {
		err := l.svcCtx.SysApiModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
			api, _ := l.svcCtx.SysApiModel.FindOne(ctx, id)
			if api == nil {
				return nil
			}
			app, _ := l.svcCtx.AppModel.FindOne(ctx, api.AppId)
			if app == nil {
				return nil
			}
			err = l.svcCtx.SysApiModel.Delete(ctx, id, db)
			if err != nil {
				return errors.Wrapf(errorx.DbError, "db api delete err:%v", err)
			}
			err = l.svcCtx.PermissionModel.Delete(ctx, api.Perms, app.Perms, db)
			if err != nil {
				return errors.Wrapf(errorx.DbError, "db permission delete err:%v", err)
			}
			res, err := l.svcCtx.CBS.DeletePolicyWithPerms(api.Perms, app.Perms)
			if err != nil || !res {
				return errors.Wrapf(errorx.NewMsg("Clear api permissions failed"), "apiPerms:%v, err:%v", api.Perms, err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.EmptyResp{}, nil
}
