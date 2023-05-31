package menulogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/treex"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuTreeLogic) GetMenuTree(in *pb.MenuQueryReq) (*pb.MenusResp, error) {
	cond := &model.SysMenuQuery{}
	_ = copier.Copiers(&cond, &in)

	menus, err := l.svcCtx.SysMenuModel.FindAll(l.ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	treeMenus := make([]model.SysMenu, 0)
	err = treex.CreateTreeFactory().ScanToTreeData(menus, &treeMenus)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("Get tree menu fail"), "err:%+v", err)
	}

	var tmr []*pb.MenuReq
	_ = copier.Copiers(&tmr, &treeMenus)

	return &pb.MenusResp{Menus: tmr}, nil
}
