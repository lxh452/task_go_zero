// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package notification

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotificationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通知列表
func NewNotificationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotificationListLogic {
	return &NotificationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotificationListLogic) NotificationList(req *types.PageReq) (resp *types.NotificationListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
