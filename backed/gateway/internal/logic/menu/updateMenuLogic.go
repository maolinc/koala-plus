package menulogic

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

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(in *pb.MenuReq) (*pb.EmptyResp, error) {
	old, _ := l.svcCtx.SysMenuModel.FindOne(l.ctx, in.Id)
	if old == nil || old.Id == 0 {
		return nil, errorx.NewMsg("Menu not exist")
	}

	var menu model.SysMenu
	_ = copier.Copiers(&menu, &in)

	menu.Perms = ""
	if menu.ParentId == menu.Id {
		menu.ParentId = 0
	}

	err := l.svcCtx.SysMenuModel.Update(l.ctx, &menu)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db update err:%v", err)
	}

	return &pb.EmptyResp{}, nil
}
