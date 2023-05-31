package rolelogic

import (
	"context"
	"encoding/json"
	"koala/gateway/internal/tools/errorx"
	utils "koala/gateway/internal/tools/treex"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageRoleLogic {
	return &PageRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageRoleLogic) PageRole(in *pb.RoleQueryReq) (*pb.RoleListResp, error) {
	cond := &model.SysRoleQuery{}
	_ = copier.Copiers(&cond, &in)
	_ = json.Unmarshal([]byte(in.SearchPlus), &cond.SearchPlus)

	roles, err := l.svcCtx.SysRoleModel.FindAll(l.ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}

	treeRole := make([]model.SysRole, 0)
	err = utils.CreateTreeFactory().ScanToTreeData(roles, &treeRole)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("Failed to obtain the tree character"), "err:%v", err)
	}

	resp := &pb.RoleListResp{Total: int64(len(treeRole))}
	_ = copier.Copiers(&resp.List, &treeRole)

	return resp, nil
}
