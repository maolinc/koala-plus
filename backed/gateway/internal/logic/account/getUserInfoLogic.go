package accountlogic

import (
	"context"
	"github.com/zeromicro/go-zero/core/fx"
	"koala/gateway/internal/tools/casbin"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v, userId:%d", err, in.UserId)
	}
	if user == nil || user.Id == 0 {
		return nil, errorx.UserNotExistError
	}
	app, _ := l.svcCtx.AppModel.FindByAppId(l.ctx, in.AppId)
	if err != nil {
		return nil, errorx.AppNotExistError
	}

	userInfoAll := &pb.GetUserInfoResp{}
	_ = copier.Copiers(&userInfoAll.User, user)

	// get permission
	if in.Role {
		roleInfo := l.GetRolePermission(user.Perms, app.Perms)
		userInfoAll.Role = roleInfo
	}
	//if in.Post {
	//	l.svcCtx.SysPostModel.FindAll()
	//}

	return userInfoAll, nil
}

func (l *GetUserInfoLogic) GetRolePermission(userPerms, appPerms string) []*pb.RolePermission {
	rps := make([]*pb.RolePermission, 0)
	roles := l.svcCtx.CBS.FindRolesForUser(userPerms, appPerms)
	permissions, _ := l.svcCtx.CBS.FindPermissionAllByUser(userPerms, appPerms)

	permsList := make([]string, len(permissions))
	permsList = append(permsList, roles...)
	for _, permission := range permissions {
		permsList = append(permsList, permission.Perms, permission.Role)
	}
	dbPermsList, _ := l.svcCtx.PermissionModel.FindByPermsList(l.ctx, permsList, appPerms)
	if dbPermsList == nil || len(dbPermsList) == 0 {
		return rps
	}
	dbPermsMap := make(map[string]*model.SysPermission)
	for _, p := range dbPermsList {
		dbPermsMap[p.Perms] = p
	}

	rpsMap := make(map[string]*pb.RolePermission)
	for _, role := range roles {
		if permsItem, ok := dbPermsMap[role]; ok {
			r := &pb.RolePermission{}
			r.RolePerms = permsItem.Perms
			r.RoleName = permsItem.Name
			r.RoleDes = permsItem.Des
			r.RoleGroup = permsItem.Group
			r.AppPerms = permsItem.AppPerms
			r.Permissions = make([]*pb.Permission, 0)
			rps = append(rps, r)
			rpsMap[role] = r
		}
	}

	if len(permissions) != 0 {
		fx.From(func(source chan<- any) {
			for _, p := range permissions {
				source <- p
			}
		}).Filter(func(item any) bool {
			_, ok := dbPermsMap[item.(*casbin.Permission).Perms]
			return ok
		}).Map(func(item any) any {
			it := item.(*casbin.Permission)
			p := dbPermsMap[it.Perms]
			pe := &pb.Permission{
				Perms:     p.Perms,
				AppPerms:  p.AppPerms,
				Name:      p.Name,
				Des:       p.Des,
				Group:     p.Group,
				RolePerms: it.Role,
			}
			return pe
		}).Group(func(item any) any {
			return item.(*pb.Permission).RolePerms
		}).ForEach(func(item any) {
			its := item.([]any)
			r := &pb.RolePermission{}
			for _, it := range its {
				p := it.(*pb.Permission)
				rp := rpsMap[p.RolePerms]
				rp.Permissions = append(r.Permissions, p)
			}
		})
	}
	return rps
}
