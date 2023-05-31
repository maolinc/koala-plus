package authoritylogic

import (
	"context"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnforceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnforceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnforceLogic {
	return &EnforceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EnforceLogic) Enforce(in *pb.PolicyReq) (*pb.BoolRep, error) {
	// todo: add your logic here and delete this line

	return &pb.BoolRep{}, nil
}
