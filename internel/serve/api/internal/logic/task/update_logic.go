package task

import (
	"context"
	"dv/internel/serve/api/internal/model"
	"github.com/jinzhu/copier"

	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.TaskUpdateRequest) (resp *types.TaskUpdateResponse, err error) {
	task := model.Task{}
	_ = copier.Copy(task, req)

	err = l.svcCtx.TaskModel.Update(task)

	return
}
