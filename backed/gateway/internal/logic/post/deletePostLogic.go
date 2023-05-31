package postlogic

import (
	"context"
	"koala/gateway/internal/tools/errorx"

	"koala/gateway/internal/svc"
	"koala/gateway/pb"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(in *pb.DeletesReq) (*pb.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return nil, errorx.NewMsg("无删除对象")
	}

	for _, id := range in.Ids {
		err := l.svcCtx.SysPostModel.Delete(l.ctx, id)
		if err != nil {
			return nil, errors.Wrapf(errorx.DbError, "db delete err:%v", err)
		}
	}

	return &pb.EmptyResp{}, nil
}
