package deptlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/treex"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageDeptLogic {
	return &PageDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageDeptLogic) PageDept(in *pb.DeptQueryReq) (*pb.DeptListResp, error) {
	cond := &model.SysDeptQuery{}
	_ = copier.Copiers(&cond, &in)

	depts, err := l.svcCtx.SysDeptModel.FindAll(l.ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	treeDept := make([]model.SysDept, 0)
	err = treex.CreateTreeFactory().ScanToTreeData(depts, &treeDept)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("获取树形部门失败"), "err:%v", err)
	}

	resp := &pb.DeptListResp{Total: int64(len(depts))}
	_ = copier.Copiers(&resp.List, &treeDept)

	return resp, nil
}
