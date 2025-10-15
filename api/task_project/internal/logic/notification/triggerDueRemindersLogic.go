// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package notification

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerDueRemindersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 触发到期提醒
func NewTriggerDueRemindersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerDueRemindersLogic {
	return &TriggerDueRemindersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerDueRemindersLogic) TriggerDueReminders(req *types.Empty) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
