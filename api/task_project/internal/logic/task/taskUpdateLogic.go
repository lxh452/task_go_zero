// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新任务
func NewTaskUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskUpdateLogic {
	return &TaskUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskUpdateLogic) TaskUpdate(req *types.TaskUpdateReq) (resp *types.Task, err error) {
	// todo: add your logic here and delete this line

	return
}
