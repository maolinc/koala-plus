package deptlogic

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

type UpdateDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDeptLogic {
	return &UpdateDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDeptLogic) UpdateDept(in *pb.DeptReq) (*pb.EmptyResp, error) {
	old, _ := l.svcCtx.SysDeptModel.FindOne(l.ctx, in.Id)
	if old == nil || old.Id == 0 {
		return nil, errorx.NewMsg("部门不存在")
	}

	var dept model.SysDept
	_ = copier.Copiers(&dept, &in)

	if dept.ParentId == dept.Id {
		dept.ParentId = 0
	}
	err := l.svcCtx.SysDeptModel.Update(l.ctx, &dept)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("部门更新失败"), "db update err:%v", err)
	}

	return &pb.EmptyResp{}, nil
}
