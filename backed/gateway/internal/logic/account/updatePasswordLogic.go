package accountlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/verifyx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(in *pb.UpdatePasswordReq) (*pb.EmptyResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db select err:%v, userId:%d", err, in.UserId)
	}

	if user == nil || user.Id == 0 {
		return nil, errors.Wrapf(errorx.UserNotExistError, "userId:%d", in.UserId)
	}

	if !verifyx.EqualsPassword(in.PasswordOld, user.Password) {
		return nil, errorx.NewMsg("旧密码错误")
	}

	pwd, err := verifyx.EncryptPassword(in.PasswordNew)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("加密失败"), "Encrypt error password:%s,err:%v", in.PasswordOld, err)
	}

	user.Password = pwd
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "db update err:%v, userId:%d", err, in.UserId)
	}

	return &pb.EmptyResp{}, nil
}
