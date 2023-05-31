package accountlogic

import (
	"context"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/verifyx"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *pb.ResetPasswordReq) (*pb.EmptyResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
	}
	if user == nil || user.Id == 0 {
		return nil, errors.Wrapf(errorx.UserNotExistError, "user not exist userId:%d", in.UserId)
	}

	pwd, err := verifyx.EncryptPassword(in.Password)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("加密失败"), "Encrypt error password:%s,err:%v", in.Password, err)
	}

	user.Password = pwd
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db update err:%v, userId:%d", err, in.UserId)
	}

	return &pb.EmptyResp{}, nil
}
