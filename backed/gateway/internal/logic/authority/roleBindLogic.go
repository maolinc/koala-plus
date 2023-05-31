package authoritylogic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/pb"
)

const (
	BindRoleRole       = "BindRoleRole"
	RemoveBindRoleRole = "RemoveBindRoleRole"
	BindRoleUser       = "BindRoleUser"
	RemoveBindRoleUser = "RemoveBindRoleUser"
	BindRoleMenu       = "BindRoleMenu"
	RemoveBindRoleMenu = "RemoveBindRoleMenu"
	BindRoleApi        = "BindRoleApi"
	RemoveBindRoleApi  = "RemoveBindRoleApi"
)

type RoleBindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleBindLogic {
	return &RoleBindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleBindLogic) RoleBind(inReq *pb.RoleBindReq) (*pb.BoolRep, error) {

	var (
		rolePerms = "xxxxxxxxxxxxxxxxxxxxxx"
		domPerms  = "xxxxxxxxxxxxxxxxxxxxxx"
	)

	bind := func(in *pb.RoleBindReq) (bool, error) {
		var resBool bool
		var err error

		if in.Role != rolePerms {
			role, _ := l.svcCtx.SysRoleModel.FindByPerms(l.ctx, in.Role)
			if role == nil {
				return false, errorx.NewPermsNotExist("Role", in.Role)
			}
			rolePerms = in.Role
		}
		if in.Dom != domPerms {
			app, _ := l.svcCtx.AppModel.FindByPerms(l.ctx, in.Dom)
			if app == nil {
				return false, errorx.NewPermsNotExist("App", in.Dom)
			}
			domPerms = in.Dom
		}

		action := in.Action
		switch action {
		case BindRoleRole:
			resBool, err = l.BindRoleRole(in)
		case RemoveBindRoleRole:
			resBool, err = l.svcCtx.CBS.RemoveBindRoleRole(in.Role, in.Obj, in.Dom)
		case BindRoleUser:
			resBool, err = l.BindRoleUser(in)
		case RemoveBindRoleUser:
			resBool, err = l.svcCtx.CBS.RemoveBindRoleUser(in.Role, in.Obj, in.Dom)
		case BindRoleMenu:
			resBool, err = l.BindRoleMenu(in)
		case RemoveBindRoleMenu:
			resBool, err = l.svcCtx.CBS.RemoveBindRoleMenu(in.Role, in.Obj, in.Dom, in.Act)
		case BindRoleApi:
			resBool, err = l.BindRoleApi(in)
		case RemoveBindRoleApi:
			resBool, err = l.svcCtx.CBS.RemoveBindRoleApi(in.Role, in.Obj, in.Dom, in.Act)
		default:
			return false, errorx.NewMsg("unknown action")
		}
		return resBool, err
	}

	var resBool bool
	var err error

	if inReq.Batch != nil && len(inReq.Batch) > 0 {
		for _, in := range inReq.Batch {
			resBool, err = bind(in)
			if err != nil {
				break
			}
		}
	} else {
		resBool, err = bind(inReq)
	}

	if err != nil {
		if _, ok := err.(*errorx.CodeError); !ok {
			err = errors.Wrapf(errorx.DbError, "err: %+v", err)
		}
		return nil, err
	}

	return &pb.BoolRep{Res: resBool}, nil
}

func (l *RoleBindLogic) BindRoleRole(in *pb.RoleBindReq) (bool, error) {
	role2, _ := l.svcCtx.SysRoleModel.FindByPerms(l.ctx, in.Obj)
	if role2 == nil {
		return false, errorx.NewPermsNotExist("Role", in.Obj)
	}

	return l.svcCtx.CBS.BindRoleRole(in.Role, in.Obj, in.Dom)
}

func (l *RoleBindLogic) BindRoleUser(in *pb.RoleBindReq) (bool, error) {
	user, _ := l.svcCtx.UserModel.FindByPerms(l.ctx, in.Role)
	if user == nil {
		return false, errorx.NewPermsNotExist("User", in.Obj)
	}

	return l.svcCtx.CBS.BindRoleUser(in.Role, in.Obj, in.Dom)
}

func (l *RoleBindLogic) BindRoleApi(in *pb.RoleBindReq) (bool, error) {
	user, _ := l.svcCtx.SysApiModel.FindByPerms(l.ctx, in.Role)
	if user == nil {
		return false, errorx.NewPermsNotExist("Api", in.Obj)
	}

	return l.svcCtx.CBS.BindRoleApi(in.Role, in.Obj, in.Dom, in.Act)
}

func (l *RoleBindLogic) BindRoleMenu(in *pb.RoleBindReq) (bool, error) {
	user, _ := l.svcCtx.SysMenuModel.FindByPerms(l.ctx, in.Role)
	if user == nil {
		return false, errorx.NewPermsNotExist("Menu", in.Obj)
	}

	return l.svcCtx.CBS.BindRoleMenu(in.Role, in.Obj, in.Dom, in.Act)
}
