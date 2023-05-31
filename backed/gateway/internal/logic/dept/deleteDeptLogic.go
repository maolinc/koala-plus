package deptlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDeptLogic {
	return &DeleteDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDeptLogic) DeleteDept(in *pb.DeletesReq) (*pb.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.NewMsg("无删除对象")
	}

	for _, id := range in.Ids {
		err := l.svcCtx.SysDeptModel.Delete(l.ctx, id)
		if err != nil {
			return nil, errors.Wrapf(errorx.DbError, "db delete err:%v", err)
		}
	}

	return &pb.EmptyResp{}, nil
}
