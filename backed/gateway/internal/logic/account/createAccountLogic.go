package accountlogic

import (
	"context"
	"github.com/maolinc/copier"
	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAccountLogic {
	return &CreateAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAccountLogic) CreateAccount(in *pb.AccountReq) (*pb.EmptyResp, error) {
	in.Id = 0
	in.CreateTime = ""
	var userReg pb.RegisterReq
	_ = copier.Copiers(&userReg, in)

	userLogic := NewRegisterUserLogic(l.ctx, l.svcCtx)
	_, err := userLogic.RegisterUser(&userReg)
	if err != nil {
		return nil, err
	}

	return &pb.EmptyResp{}, nil
}
