// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新任务明细/进度
func NewLogUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogUpdateLogic {
	return &LogUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogUpdateLogic) LogUpdate(req *types.LogUpdateReq) (resp *types.UserTaskLog, err error) {
	// todo: add your logic here and delete this line

	return
}
