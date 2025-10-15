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
	// 1. 查询所有任务日志
	logs, err := l.svcCtx.UserTaskLogModel.FindAll(l.ctx, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 2. 查找即将到期的任务
	now := time.Now()
	reminderThreshold := now.Add(24 * time.Hour) // 24小时内到期
	
	var dueTasks []*core.UserTaskLog
	for _, log := range logs {
		if !log.DueDate.IsZero() && log.DueDate.Before(reminderThreshold) && log.DueDate.After(now) {
			dueTasks = append(dueTasks, log)
		}
	}

	// 3. 发送到期提醒邮件
	reminderCount := 0
	for _, log := range dueTasks {
		// 获取任务信息
		task, err := l.svcCtx.TaskModel.FindOne(l.ctx, log.TaskId)
		if err != nil {
			continue
		}

		// 获取用户信息
		user, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, log.UserId)
		if err != nil {
			continue
		}

		// 发送提醒邮件给负责人
		for _, userId := range task.ResponsibleUserIds {
			responsibleUser, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
			if err != nil {
				continue
			}
			
			l.svcCtx.Mailer.SendEmail(l.ctx, responsibleUser.Email, "任务即将到期提醒", 
				fmt.Sprintf("任务 %s 即将在 %s 到期，请及时处理。", task.Title, log.DueDate.Format("2006-01-02 15:04:05")))
			reminderCount++
		}

		// 发送提醒邮件给节点用户
		for _, userId := range task.NodeUserIds {
			nodeUser, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
			if err != nil {
				continue
			}
			
			l.svcCtx.Mailer.SendEmail(l.ctx, nodeUser.Email, "任务节点即将到期提醒", 
				fmt.Sprintf("任务节点 %s 即将在 %s 到期，请及时处理。", task.Title, log.DueDate.Format("2006-01-02 15:04:05")))
			reminderCount++
		}
	}

	// 4. 发布到期提醒事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.reminder.due", map[string]interface{}{
		"due_tasks_count": len(dueTasks),
		"reminders_sent":  reminderCount,
		"triggered_at":    now,
	})

	return &types.ActionResp{Success: true, Message: fmt.Sprintf("已发送 %d 个到期提醒", reminderCount)}, nil
}
