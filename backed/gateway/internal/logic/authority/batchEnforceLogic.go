package authoritylogic

import (
	"context"
	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchEnforceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchEnforceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchEnforceLogic {
	return &BatchEnforceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchEnforceLogic) BatchEnforce(in *pb.RulesReq) (*pb.BoolRep, error) {
	// todo: add your logic here and delete this line

	return &pb.BoolRep{}, nil
}
