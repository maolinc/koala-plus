package applicationlogic

import (
	"context"
	"github.com/pkg/errors"
	"koala/gateway/internal/tools/casbin"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/utilx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSysApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSysApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSysApplicationLogic {
	return &CreateSysApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSysApplicationLogic) CreateSysApplication(in *pb.CreateSysApplicationReq) (*pb.CreateSysApplicationResp, error) {
	if in.Name == "" || in.OrganizationId == 0 {
		return nil, errorx.NewMsg("name, organization required")
	}
	app := &model.SysApplication{}
	app.Name = in.Name
	app.Des = in.Des
	// 生成appId
	app.AppId = utilx.GenerateRandomString(5)
	app.Perms = utilx.GetNewPerms(app.AppId, casbin.CasDom)

	err := l.svcCtx.AppModel.Insert(l.ctx, app)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db app insert error: %+v", err)
	}

	return &pb.CreateSysApplicationResp{AppId: app.AppId, Perms: app.Perms}, nil
}
