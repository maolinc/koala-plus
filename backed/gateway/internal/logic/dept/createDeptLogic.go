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

type CreateDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDeptLogic {
	return &CreateDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDeptLogic) CreateDept(in *pb.DeptReq) (*pb.EmptyResp, error) {
	var dept model.SysDept
	_ = copier.Copiers(&dept, in)
	dept.Id = 0

	if dept.ParentId == dept.Id {
		dept.ParentId = 0
	}
	if dept.ParentId != 0 {
		pDept, _ := l.svcCtx.SysDeptModel.FindOne(l.ctx, dept.ParentId)
		if pDept.Id == 0 {
			return nil, errors.Wrapf(errorx.NewMsg("上级部门不存在"), "parentId:%d", dept.ParentId)
		}
	}

	err := l.svcCtx.SysDeptModel.Insert(l.ctx, &dept)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db insert err:%v", err)
	}

	return &pb.EmptyResp{}, nil
}
