package accountlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAccountLogic {
	return &UpdateAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAccountLogic) UpdateAccount(in *pb.AccountReq) (*pb.EmptyResp, error) {
	in.UserName = ""
	in.Password = ""
	var user model.SysUser
	_ = copier.Copiers(&user, in)

	err := l.svcCtx.UserModel.Update(l.ctx, &user)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("账号更新失败"), "db err:%v, req:%v", err, user)
	}

	// todo 岗位

	return &pb.EmptyResp{}, nil
}
