package apilogic

import (
	"context"
	"github.com/maolinc/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/casbin"
	"koala/gateway/internal/tools/constant"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/stringx"
	"koala/gateway/internal/tools/utilx"
	"koala/gateway/pb"
	"koala/model"
)

type CreateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateApiLogic) CreateApi(in *pb.ApiReq) (*pb.EmptyResp, error) {
	var api model.SysApi
	_ = copier.Copiers(&api, in)
	api.Id = 0
	app, _ := l.svcCtx.AppModel.FindOne(l.ctx, api.AppId)
	if app == nil {
		return nil, errorx.AppNotExistError
	}

	api.Path = stringx.HasPrefixAndJoin(api.Path, "/")
	m, ok := stringx.FixHttpMethod(api.Method)
	if !ok {
		return nil, errorx.NewMsg("Http method error")
	}
	api.Method = m

	old, _ := l.svcCtx.SysApiModel.FindByPath(l.ctx, api.Path, api.Method)
	if old != nil {
		return nil, errorx.NewMsg("Api is exist")
	}
	api.Perms = utilx.GetNewPerms(api.Method+api.Path, casbin.CasApi)

	err := l.svcCtx.SysApiModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
		err = l.svcCtx.SysApiModel.Insert(ctx, &api, db)
		if err != nil {
			return errors.Wrapf(errorx.SystemError, "db api insert err:%v", err)
		}

		err = l.svcCtx.PermissionModel.Insert(ctx, &model.SysPermission{
			Perms:    api.Perms,
			Name:     api.Name,
			Des:      api.Des,
			Status:   "1",
			AppPerms: app.Perms,
			Group:    constant.PERMS_GROUP_API,
			Table:    api.TableName(),
			CreateBy: 0,
		}, db)
		if err != nil {
			return errors.Wrapf(errorx.SystemError, "db permission insert err:%v", err)
		}
		return nil
	})

	return &pb.EmptyResp{}, err
}
