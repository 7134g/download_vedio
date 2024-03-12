package task

import (
	"context"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"
	"dv/internel/serve/api/internal/util/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunLogic {
	return &RunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunLogic) Run(req *types.TaskRunRequest) (resp *types.TaskRunResponse, err error) {
	//return &types.TaskRunResponse{Message: "停止成功"}, err
	if req.Stop {
		l.svcCtx.TaskControl.Stop()
		return &types.TaskRunResponse{Message: "停止成功"}, err
	}

	task := make([]model.Task, 0)
	_db := l.svcCtx.TaskModel.DB
	if len(req.IDS) != 0 {
		_db = _db.Where("id IN ?", req.IDS)
	}
	_db = _db.Where("status != ?", model.StatusSuccess)
	_db.Find(&task)

	if len(task) == 0 {
		return &types.TaskRunResponse{Message: "无任务"}, err
	}
	//if l.svcCtx.TaskControl.GetStatus() {
	//	return &types.TaskRunResponse{Message: "正在执行中"}, err
	//}
	go l.svcCtx.TaskControl.Run(task)

	return &types.TaskRunResponse{Message: "开始运行..."}, nil
}
