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
	// 1. 验证用户是否存在
	user, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.User_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在")
		}
		return nil, err
	}

	// 2. 验证通知是否存在（使用UserTaskLog模拟）
	log, err := l.svcCtx.UserTaskLogModel.FindOne(l.ctx, req.Notification_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("通知不存在")
		}
		return nil, err
	}

	// 3. 验证用户是否有权限确认该通知
	if log.UserId != req.User_id {
		return nil, logx.Errorf("用户无权限确认该通知")
	}

	// 4. 删除Redis缓存（模拟确认）
	logKey := "notification:" + string(rune(log.Id))
	l.svcCtx.Redis.Del(logKey)

	// 5. 发布通知确认事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "notification.ack", map[string]interface{}{
		"notification_id": req.Notification_id,
		"user_id":         req.User_id,
		"acknowledged_at": time.Now(),
	})

	return &types.ActionResp{Success: true, Message: "通知确认成功"}, nil
}
