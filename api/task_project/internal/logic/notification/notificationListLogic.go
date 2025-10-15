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
	// 1. 查询任务日志（模拟通知）
	logs, err := l.svcCtx.UserTaskLogModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为通知格式
	var notifications []types.Notification
	for _, log := range logs {
		notifications = append(notifications, types.Notification{
			Id:          log.Id,
			User_id:     log.UserId,
			Task_id:     log.TaskId,
			Title:       l.getNotificationTitle(log),
			Content:     l.getNotificationContent(log),
			Category:    l.getNotificationCategory(log),
			Is_read:     false, // 简化实现
			Created_at:  log.CreatedAt.String(),
			Updated_at:  log.UpdatedAt.String(),
		})
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "notification.list_queried", map[string]interface{}{
		"page":         req.Page,
		"size":         req.Size,
		"count":        len(notifications),
		"queried_at":   time.Now(),
	})

	return &types.NotificationListResp{
		List:  notifications,
		Total: int64(len(notifications)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// 获取通知标题
func (l *NotificationListLogic) getNotificationTitle(log *core.UserTaskLog) string {
	if log.Progress == 100 {
		return "任务完成"
	} else if log.Progress == -1 {
		return "任务交接被拒绝"
	} else if log.HandoverId != 0 {
		return "任务交接"
	} else {
		return "任务更新"
	}
}

// 获取通知内容
func (l *NotificationListLogic) getNotificationContent(log *core.UserTaskLog) string {
	if log.Progress == 100 {
		return "您的任务已完成"
	} else if log.Progress == -1 {
		return "任务交接被拒绝"
	} else if log.HandoverId != 0 {
		return "有新的任务交接"
	} else {
		return "任务进度已更新"
	}
}

// 获取通知分类
func (l *NotificationListLogic) getNotificationCategory(log *core.UserTaskLog) string {
	if log.Progress == 100 {
		return "task_completed"
	} else if log.Progress == -1 {
		return "handover_rejected"
	} else if log.HandoverId != 0 {
		return "handover"
	} else {
		return "task_update"
	}
}
