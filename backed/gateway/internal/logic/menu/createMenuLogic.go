package menulogic

import (
	"context"
	"gorm.io/gorm"
	"koala/gateway/internal/tools/casbin"
	"koala/gateway/internal/tools/constant"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/utilx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	MENU_TYPE_M = "M"
	MENU_TYPE_C = "C"
	MENU_TYPE_F = "F"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(in *pb.MenuReq) (*pb.EmptyResp, error) {
	var menu model.SysMenu
	_ = copier.Copiers(&menu, in)
	menu.Id = 0

	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, menu.AppId)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	if menu.ParentId == menu.Id {
		menu.ParentId = 0
	}
	if menu.ParentId != 0 {
		pMenu, _ := l.svcCtx.SysMenuModel.FindOne(l.ctx, menu.ParentId)
		if pMenu.Id == 0 {
			return nil, errorx.NewMsg("Parent menu not exist")
		}
	}

	if menu.MenuType == MENU_TYPE_C || menu.MenuType == MENU_TYPE_F {
		if menu.Perms == "" {
			return nil, errorx.PermsNotEmptyError
		}
		menu.Perms = utilx.GetNewPerms(menu.Perms, casbin.CasMenu)
		old, _ := l.svcCtx.SysMenuModel.FindByPerms(l.ctx, menu.Perms)
		if old != nil && old.Id != 0 {
			return nil, errorx.NewPermsExist("menu", menu.Perms)
		}
	}

	err := l.svcCtx.SysMenuModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
		err = l.svcCtx.SysMenuModel.Insert(ctx, &menu)
		if err != nil {
			return errors.Wrapf(errorx.DbError, "db menu insert err:%v", err)
		}
		err = l.svcCtx.PermissionModel.Insert(ctx, &model.SysPermission{
			Perms:    menu.Perms,
			Name:     menu.MenuName,
			Des:      menu.MenuName,
			Status:   "1",
			AppPerms: app.Perms,
			Group:    constant.PERMS_GROUP_MENU,
			Table:    menu.TableName(),
			CreateBy: 0,
		})
		if err != nil {
			return errors.Wrapf(errorx.DbError, "db permission insert err:%v", err)
		}
		return nil
	})

	return &pb.EmptyResp{}, err
}
