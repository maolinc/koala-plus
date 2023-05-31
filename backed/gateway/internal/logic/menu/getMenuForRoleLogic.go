package menulogic

import (
	"context"
	"github.com/pkg/errors"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuForRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuForRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuForRoleLogic {
	return &GetMenuForRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuForRoleLogic) GetMenuForRole(in *pb.PermsReq) (*pb.PermsResp, error) {
	menu, err := l.svcCtx.SysMenuModel.FindByPerms(l.ctx, in.Perms)
	if err != nil || menu == nil {
		return nil, errorx.NewMsg("Menu not exist")
	}
	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, menu.AppId)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	menuPermsList, err := l.svcCtx.CBS.FindMenusForRole(menu.Perms, app.Perms)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("Get menu permission fail"), "role perms:%s, err:%v", in.Perms, err)
	}

	return &pb.PermsResp{PermsList: menuPermsList}, nil
}
