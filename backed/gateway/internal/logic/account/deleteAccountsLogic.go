package accountlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAccountsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAccountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAccountsLogic {
	return &DeleteAccountsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAccountsLogic) DeleteAccounts(in *pb.DeletesReq) (*pb.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.NewMsg("无删除对象")
	}

	for _, id := range in.Ids {
		user, err := l.svcCtx.UserModel.FindOne(l.ctx, id)
		if err != nil {
			return nil, errors.Wrapf(errorx.DbError, "db select err:%v", err)
		}
		if user == nil {
			continue
		}
		_, err = l.svcCtx.CBS.CB.DeleteRole(user.Perms)
		if err != nil {
			return nil, errors.Wrapf(errorx.NewMsgF("Clear user %s role permission fail", user.UserName),
				"userId:%v, perms:%v, err:%v", id, user.Perms, err)
		}
		err = l.svcCtx.UserModel.Delete(l.ctx, id)
		if err != nil {
			return nil, errors.Wrapf(errorx.DbError, "db delete err:%v", err)
		}
	}

	return &pb.EmptyResp{}, nil
}
