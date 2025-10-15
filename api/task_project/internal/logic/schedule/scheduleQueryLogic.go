// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schedule

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScheduleQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询任务排期
func NewScheduleQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleQueryLogic {
	return &ScheduleQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleQueryLogic) ScheduleQuery(req *types.ScheduleQueryReq) (resp *types.LogListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
