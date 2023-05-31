package apilogic

import (
	"context"
	"gorm.io/gorm"
	"koala/gateway/internal/tools/casbin"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/stringx"
	"koala/gateway/internal/tools/utilx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApiLogic) UpdateApi(in *pb.ApiReq) (*pb.EmptyResp, error) {
	old, _ := l.svcCtx.SysApiModel.FindOne(l.ctx, in.Id)
	if old == nil || old.Id == 0 {
		return nil, errorx.NewMsg("Api does not exist")
	}

	var api model.SysApi
	_ = copier.Copiers(&api, &in)
	api.AppId = 0

	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, old.AppId)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	if api.Method != "" {
		m, ok := stringx.FixHttpMethod(api.Method)
		if !ok {
			return nil, errorx.NewMsg("Http method error")
		}
		api.Method = m
	} else {
		api.Method = old.Method
	}

	if api.Path != "" {
		api.Path = stringx.HasPrefixAndJoin(api.Path, "/")
	} else {
		api.Path = old.Path
	}

	api.Perms = utilx.GetNewPerms(api.Method+api.Path, casbin.CasApi)

	err := l.svcCtx.SysApiModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
		err = l.svcCtx.SysApiModel.Update(l.ctx, &api)
		if err != nil {
			return errors.Wrapf(errorx.NewMsg("Api update fail"), "db update err:%v", err)
		}
		err = l.svcCtx.PermissionModel.UpdateByPerms(l.ctx, old.Perms, app.Perms, &model.SysPermission{
			Perms: api.Perms,
			Name:  api.Name,
			Des:   api.Des,
		})
		if err != nil {
			return errors.Wrapf(errorx.NewMsg("Api update permission fail"), "db update err:%v", err)
		}

		if old.Perms != api.Perms {
			res := l.svcCtx.CBS.UpdatePolicyWithPerms(old.Perms, api.Perms, app.Perms)
			if !res {
				return errorx.NewMsg("Permission move fail, you can try delete, then create api")
			}
		}
		return nil
	})

	return &pb.EmptyResp{}, err
}
