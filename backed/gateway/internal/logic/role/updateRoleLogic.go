package rolelogic

import (
	"context"
	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"koala/gateway/internal/tools/collectx"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *pb.RoleReq) (*pb.EmptyResp, error) {
	old, _ := l.svcCtx.SysRoleModel.FindOne(l.ctx, in.Id)
	if old == nil || old.Id == 0 {
		return nil, errorx.NewMsg("Role is not exist")
	}

	var role model.SysRole
	_ = copier.Copiers(&role, &in)
	if role.ParentId == role.Id {
		role.ParentId = 0
	}
	role.AppId = 0  // app不允许更更新
	role.Perms = "" // 权限字符不允许更新
	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, old.AppId)
	if app == nil {
		return nil, errorx.AppNotExistError
	}
	dom := app.Perms

	err := l.svcCtx.SysRoleModel.Update(l.ctx, &role)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("Role update fail"), "db update err:%v", err)
	}

	existMenuPerms, err := l.svcCtx.CBS.FindMenusForRole(in.Perms, dom)

	_, d1, d2 := collectx.JDiffString(existMenuPerms, in.MenuPerms)
	if len(d1) > 0 {
		for _, menuPerms := range d1 {
			_, _ = l.svcCtx.CBS.RemoveBindRoleMenu(old.Perms, menuPerms, dom, "*")
		}
	}
	if len(d2) > 0 {
		menus, _ := l.svcCtx.SysMenuModel.FindByPermsList(l.ctx, d2)
		for _, menu := range menus {
			_, _ = l.svcCtx.CBS.BindRoleMenu(old.Perms, menu.Perms, dom, "*")
		}
	}

	return &pb.EmptyResp{}, nil
}
