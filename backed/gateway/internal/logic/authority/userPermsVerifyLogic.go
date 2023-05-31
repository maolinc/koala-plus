package authoritylogic

import (
	"context"
	"github.com/pkg/errors"
	"koala/gateway/internal/svc"
	"koala/gateway/internal/tools/errorx"
	"koala/gateway/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserPermsVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserPermsVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserPermsVerifyLogic {
	return &UserPermsVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserPermsVerifyLogic) UserPermsVerify(in *pb.UserPermsVerifyReq) (*pb.BoolRep, error) {

	var resBool bool
	var err error

	if in.Batch != nil && len(in.Batch) > 0 {
		for _, r := range in.Batch {
			resBool, err = l.svcCtx.CBS.CB.Enforce(r.User, r.Dom, r.Obj, r.Act)
			if err != nil || !resBool {
				break
			}
		}
	} else {
		resBool, err = l.svcCtx.CBS.CB.Enforce(in.User, in.Dom, in.Obj, in.Act)
	}

	if err != nil {
		return &pb.BoolRep{Res: false}, errors.Wrapf(errorx.SystemError, "err: %+v", err)
	}

	return &pb.BoolRep{Res: resBool}, nil
}
