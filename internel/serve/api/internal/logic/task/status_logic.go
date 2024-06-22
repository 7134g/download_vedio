package task

import (
	"context"
	"dv/internel/serve/api/internal/svc/dtask"

	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusLogic {
	return &StatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatusLogic) Status(req *types.TaskStatusRequest) (resp *types.TaskStatusResponse, err error) {
	resp = &types.TaskStatusResponse{
		Status:   dtask.Control.CheckStatus(),
		WebProxy: l.svcCtx.Config.WebProxy,
	}

	return
}
