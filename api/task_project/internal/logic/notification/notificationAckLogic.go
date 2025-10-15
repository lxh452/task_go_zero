// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package notification

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationAckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通知已读确认
func NewNotificationAckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationAckLogic {
	return &NotificationAckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotificationAckLogic) NotificationAck(req *types.NotificationAckReq) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
