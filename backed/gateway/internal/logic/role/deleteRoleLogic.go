package rolelogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"

	"gorm.io/gorm"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteRoleLogic) DeleteRole(in *pb.DeletesReq) (*pb.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return nil, nil
	}

	roles, _ := l.svcCtx.SysRoleModel.FindByIds(l.ctx, in.Ids)

	err := l.svcCtx.SysRoleModel.Trans(l.ctx, func(ctx context.Context, db *gorm.DB) (err error) {
		for _, role := range roles {
			err = l.svcCtx.SysRoleModel.Delete(l.ctx, role.Id)
			if err != nil {
				return errors.Wrapf(errorx.DbError, "db delete err:%v", err)
			}
			res, err := l.svcCtx.CBS.CB.DeleteRole(role.Perms)
			if err != nil || !res {
				return errors.Wrapf(errorx.NewMsg("Clear role permissions failed"), "role perms:%v, err:%v", role.Perms, err)
			}
		}
		return nil
	})

	return &pb.EmptyResp{}, err
}
