package menulogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"

	"gorm.io/gorm"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenusLogic {
	return &DeleteMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenusLogic) DeleteMenus(in *pb.DeletesReq) (*pb.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return nil, nil
	}

	for _, id := range in.Ids {
		err := l.svcCtx.SysMenuModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
			menu, _ := l.svcCtx.SysMenuModel.FindOne(ctx, id)
			if menu == nil {
				return nil
			}
			app, _ := l.svcCtx.AppModel.FindOne(ctx, menu.AppId)
			if app == nil {
				return errorx.AppNotExistError
			}

			err = l.svcCtx.SysMenuModel.Delete(ctx, id, db)
			if err != nil {
				return errors.Wrapf(errorx.DbError, "db menu delete err:%v", err)
			}

			err = l.svcCtx.PermissionModel.Delete(ctx, menu.Perms, app.Perms, db)
			if err != nil {
				return errors.Wrapf(errorx.DbError, "db permission delete err:%v", err)
			}
			// 删除casbin的菜单权限
			res, err := l.svcCtx.CBS.DeletePolicyWithPerms(menu.Perms, app.Perms)
			if err != nil || !res {
				return errors.Wrapf(errorx.NewMsg("Clear menu permissions failed"), "menuPerms:%v, err:%v", menu.Perms, err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return &pb.EmptyResp{}, nil
}
