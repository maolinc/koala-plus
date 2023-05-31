package rolelogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRoleLogic) CreateRole(in *pb.RoleReq) (*pb.EmptyResp, error) {
	var role model.SysRole
	_ = copier.Copiers(&role, in)
	role.Id = 0

	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, role.AppId)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	if role.ParentId == role.Id {
		role.ParentId = 0
	}

	var pRole *model.SysRole
	if role.ParentId != 0 {
		pRole, _ = l.svcCtx.SysRoleModel.FindOne(l.ctx, role.ParentId)
		if pRole == nil || pRole.Id == 0 {
			return nil, errors.Wrapf(errorx.NewMsg("Parent role not exist"), "parentId:%v", &role.ParentId)
		}
	}

	old, _ := l.svcCtx.SysRoleModel.FindByPerms(l.ctx, role.Perms)
	if old != nil && old.Id != 0 {
		return nil, errorx.NewPermsExist("role", role.Perms)
	}

	err := l.svcCtx.SysRoleModel.Insert(l.ctx, &role)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db insert err:%v", err)
	}

	in.Id = role.Id

	if len(in.MenuPerms) > 0 {
		menus, _ := l.svcCtx.SysMenuModel.FindByPermsList(l.ctx, in.MenuPerms)
		l.bindMenus(menus, app.Perms, role.Perms)
	}

	if role.ParentId != 0 {
		res, _ := l.svcCtx.CBS.BindRoleRole(pRole.Perms, role.Perms, app.Perms)
		if !res {
			return nil, errorx.NewCodeMsgF(errorx.Perms, "Permission inheritance failed")
		}
	}

	return &pb.EmptyResp{}, nil
}

func (l *CreateRoleLogic) bindMenus(menus []model.SysMenu, role, dom string) {
	for _, menu := range menus {
		_, _ = l.svcCtx.CBS.BindRoleMenu(role, menu.Perms, dom, "*")
	}
}
