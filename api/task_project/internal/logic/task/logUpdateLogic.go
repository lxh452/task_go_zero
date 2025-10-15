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
	// 1. 查找任务日志
	userTaskLog, err := l.svcCtx.UserTaskLogModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务日志不存在")
		}
		return nil, err
	}

	// 2. 验证协同用户是否存在（如果有变化）
	if req.Collaborator_id != 0 && req.Collaborator_id != userTaskLog.CollaboratorId {
		collaborator, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.Collaborator_id)
		if err != nil {
			if err == core.ErrNotFound {
				return nil, logx.Errorf("协同用户不存在")
			}
			return nil, err
		}
		if collaborator.Status != 1 {
			return nil, logx.Errorf("协同用户不在职")
		}
	}

	// 3. 更新任务日志信息
	now := time.Now()
	oldProgress := userTaskLog.Progress
	
	if req.Progress != 0 {
		userTaskLog.Progress = req.Progress
	}
	if req.Collaborator_id != 0 {
		userTaskLog.CollaboratorId = req.Collaborator_id
	}
	if !req.Due_date.IsZero() {
		userTaskLog.DueDate = req.Due_date
	}
	if req.Handover_id != 0 {
		userTaskLog.HandoverId = req.Handover_id
	}
	userTaskLog.UpdatedAt = now

	err = l.svcCtx.UserTaskLogModel.Update(l.ctx, userTaskLog)
	if err != nil {
		return nil, err
	}

	// 4. 更新Redis缓存
	logKey := "task_log:" + string(rune(userTaskLog.Id))
	l.svcCtx.Redis.Setex(logKey, 3600, userTaskLog) // 缓存1小时

	// 5. 发布任务日志更新事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.log.updated", map[string]interface{}{
		"log_id":        userTaskLog.Id,
		"user_id":       userTaskLog.UserId,
		"task_id":       userTaskLog.TaskId,
		"progress":      userTaskLog.Progress,
		"updated_at":    now,
	})

	// 6. 如果进度有变化，发送通知邮件
	if oldProgress != userTaskLog.Progress {
		// 获取任务信息
		task, _ := l.svcCtx.TaskModel.FindOne(l.ctx, userTaskLog.TaskId)
		if task != nil {
			for _, userId := range task.ResponsibleUserIds {
				if userId != userTaskLog.UserId {
					user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
					l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "任务进度更新", 
						fmt.Sprintf("任务 %s 进度已更新为 %d%%", task.Title, userTaskLog.Progress))
				}
			}
		}
	}

	return &types.UserTaskLog{
		Id:              userTaskLog.Id,
		User_id:         userTaskLog.UserId,
		Task_id:         userTaskLog.TaskId,
		Node_id:         userTaskLog.NodeId,
		Collaborator_id: userTaskLog.CollaboratorId,
		Progress:        userTaskLog.Progress,
		Due_date:        userTaskLog.DueDate.String(),
		Handover_id:     userTaskLog.HandoverId,
		Created_at:      userTaskLog.CreatedAt.String(),
		Updated_at:      userTaskLog.UpdatedAt.String(),
	}, nil
}
