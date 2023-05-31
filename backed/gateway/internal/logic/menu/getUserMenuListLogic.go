package menulogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	utils "koala/gateway/internal/tools/treex"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMenuListLogic {
	return &GetUserMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserMenuListLogic) GetUserMenuList(in *pb.UserReq) (*pb.UserMenusResp, error) {
	user, _ := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if user == nil {
		return nil, errorx.UserNotExistError
	}
	menuPerms, err := l.svcCtx.CBS.FindMenusForUser(user.Perms, in.Dom)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("Get user menu list fail"), "userId:%v, dom:%v, err:%+v", in.UserId, in.Dom, err)
	}

	menus, err := l.svcCtx.SysMenuModel.FindByPermsList(l.ctx, menuPerms)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	treeMenus := make([]model.SysMenu, 0)
	err = utils.CreateTreeFactory().ScanToTreeData(menus, &treeMenus)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("Get tree menu fail"), "err:%+v", err)
	}

	mR := make([]*pb.MenusListResp, 0)
	tree(treeMenus, &mR)

	return &pb.UserMenusResp{Menus: mR}, nil
}

func tree(from []model.SysMenu, to *[]*pb.MenusListResp) {
	for _, menu := range from {
		me := pb.MenusListResp_MenuMeta{
			Title:             menu.MenuName,
			CurrentActiveMenu: menu.CurrentActiveMenu,
			HideMenu:          menu.Visible == "1",
			Icon:              menu.Icon,
			OrderNo:           menu.OrderNo,
			IgnoreKeepAlive:   menu.IsCache == 0,
			FrameSrc:          "",
			IgnoreRoute:       false,
		}
		mR := pb.MenusListResp{
			Id:        menu.Id,
			Name:      menu.MenuName,
			ParentId:  menu.ParentId,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Perms:     menu.Perms,
			Meta:      &me,
		}
		if len(menu.Children) > 0 {
			mRR := make([]*pb.MenusListResp, 0)
			tree(menu.Children, &mRR)
			mR.Children = mRR
		}
		*to = append(*to, &mR)
	}
}
