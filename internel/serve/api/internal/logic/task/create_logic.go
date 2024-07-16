package task

import (
	"context"
	"dv/internel/serve/api/internal/model"
	"github.com/jinzhu/copier"

	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.TaskCreateRequest) (resp *types.TaskCreateResponse, err error) {
	task := model.Task{}
	_ = copier.Copy(task, req)
	err = l.svcCtx.TaskModel.Create(&task)

	return nil, err
}
