// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增任务明细/日报
func NewLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogCreateLogic {
	return &LogCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogCreateLogic) LogCreate(req *types.LogCreateReq) (resp *types.UserTaskLog, err error) {
	// todo: add your logic here and delete this line

	return
}
