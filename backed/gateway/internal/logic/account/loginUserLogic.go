package accountlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/internal/tools/verifyx"
	"koala/model"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	LOGIN_PHONE    = "phone"
	LOGIN_USRENAME = "username"
	LOGIN_ACCOUNT  = "account"

	account_lock = "1"
)

type LoginUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginUserLogic {
	return &LoginUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginUserLogic) LoginUser(in *pb.LoginReq) (*pb.LoginResp, error) {
	var (
		user     *model.SysUser
		err      error
		authKey  = in.AuthKey
		authType = in.AuthType
	)

	if authType == LOGIN_PHONE {
		user, err = l.loginByPhone(authKey, in.Password)
	} else if authType == LOGIN_USRENAME {
		user, err = l.loginByUsername(authKey, in.Password)
	} else if authType == LOGIN_ACCOUNT {
		if verifyx.VerifyMobileFormat(authKey) {
			user, err = l.loginByPhone(authKey, in.Password)
		} else {
			user, err = l.loginByUsername(authKey, in.Password)
		}
	} else {
		return nil, errors.Wrapf(errorx.NewMsg("登录类型值有误"), "Login authType error authType:%v", in.AuthKey)
	}

	if err != nil {
		return nil, err
	}

	if user.Status == account_lock {
		return nil, errorx.NewMsg("账号被锁定，请联系管理员")
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{UserId: user.Id})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

func (l *LoginUserLogic) loginByPhone(phone, password string) (*model.SysUser, error) {

	if !verifyx.VerifyMobileFormat(phone) {
		return nil, errors.Wrapf(errorx.NewMsg("手机号格式有误"), "phone:%v", phone)
	}

	user, err := l.svcCtx.UserModel.FindByPhone(l.ctx, phone)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "Login db select err:%v, phone:%v", err, phone)
	}
	if user == nil || user.Id == 0 {
		return nil, errors.Wrapf(errorx.UserNotExistError, "Login phone:%v", phone)

	}

	if !verifyx.EqualsPassword(password, user.Password) {
		return nil, errorx.NewMsg("密码错误")
	}

	return user, nil
}

func (l *LoginUserLogic) loginByUsername(userName, password string) (*model.SysUser, error) {

	if userName == "" {
		return nil, errors.Wrapf(errorx.NewMsg("用户名为空"), "")
	}

	user, err := l.svcCtx.UserModel.FindByUserName(l.ctx, userName)
	if err != nil {
		return nil, errors.Wrapf(errorx.DbError, "Login db select err:%v, userName:%v", err, userName)
	}
	if user == nil || user.Id == 0 {
		return nil, errors.Wrapf(errorx.UserNotExistError, "Login userName:%v", userName)
	}
	if !verifyx.EqualsPassword(password, user.Password) {
		return nil, errorx.NewMsg("密码错误")
	}

	return user, nil
}
