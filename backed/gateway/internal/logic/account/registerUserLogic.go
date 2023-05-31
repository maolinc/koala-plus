package accountlogic

import (
	"context"
	"fmt"
	"koala/gateway/internal/tools/casbin"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/verifyx"
	"koala/model"

	"github.com/maolinc/copier"
	"github.com/pkg/errors"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterUserLogic {
	return &RegisterUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterUserLogic) RegisterUser(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	if !verifyx.VerifyMobileFormat(in.Phone) || in.Password == "" {
		return nil, errors.Wrapf(errorx.NewMsg("手机号或密码格式有误"), "Register userInfo errorr phone:%v", in.Phone)
	}
	if in.UserName == "" {
		return nil, errors.Wrapf(errorx.NewMsg("用户名为空"), "Register userInfo errorr phone:%v", in.Phone)
	}

	user, err := l.svcCtx.UserModel.FindByUserName(l.ctx, in.UserName)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "Register db select err:%v,username:%s", err, in.UserName)
	}
	if user != nil && user.Id != 0 {
		return nil, errors.Wrapf(errorx.NewMsg("用户名已存在"), "Register user exist err:%v,phone:%s", err, in.Phone)
	}

	user, err = l.svcCtx.UserModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "Register db select err:%v,phone:%s", err, in.Phone)
	}
	if user != nil && user.Id != 0 {
		return nil, errors.Wrapf(errorx.NewMsg("手机号已存在"), "Register user exist err:%v,phone:%s", err, in.Phone)
	}

	var registerUser model.SysUser

	_ = copier.Copiers(&registerUser, in)
	pwd, err := verifyx.EncryptPassword(in.Password)
	if err != nil {
		return nil, errors.Wrapf(errorx.NewMsg("加密失败"), "Encrypt error password:%s,err:%v", in.Password, err)
	}
	registerUser.Password = pwd

	err = l.svcCtx.UserModel.Insert(l.ctx, &registerUser)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "Register db insert err:%v,user:%v", err, registerUser)
	}
	registerUser.Perms = fmt.Sprintf("%s%d", casbin.CasUser, registerUser.Id)
	err = l.svcCtx.UserModel.Update(l.ctx, &registerUser)
	if err != nil {
		_ = l.svcCtx.UserModel.ForceDelete(l.ctx, registerUser.Id)
		return nil, errors.Wrapf(errorx.DbError, "Register db insert err:%v, perms:%v", err, registerUser.Perms)
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.generateToken(&pb.GenerateTokenReq{UserId: registerUser.Id})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
		Id:           registerUser.Id,
	}, nil
}
