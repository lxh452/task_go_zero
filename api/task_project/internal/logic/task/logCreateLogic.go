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
	// 1. 验证用户是否存在
	user, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.User_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在")
		}
		return nil, err
	}
	if user.Status != 1 {
		return nil, logx.Errorf("用户不在职")
	}

	// 2. 验证任务是否存在
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Task_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务不存在")
		}
		return nil, err
	}

	// 3. 验证用户是否有权限操作该任务
	hasPermission := false
	for _, userId := range task.ResponsibleUserIds {
		if userId == req.User_id {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		for _, userId := range task.NodeUserIds {
			if userId == req.User_id {
				hasPermission = true
				break
			}
		}
	}
	if !hasPermission {
		return nil, logx.Errorf("用户无权限操作该任务")
	}

	// 4. 验证协同用户是否存在（如果有）
	if req.Collaborator_id != 0 {
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

	// 5. 创建任务日志
	now := time.Now()
	userTaskLog := &core.UserTaskLog{
		UserId:         req.User_id,
		TaskId:         req.Task_id,
		NodeId:         req.Node_id,
		CollaboratorId: req.Collaborator_id,
		Progress:       req.Progress,
		DueDate:        req.Due_date,
		HandoverId:     req.Handover_id,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = l.svcCtx.UserTaskLogModel.Insert(l.ctx, userTaskLog)
	if err != nil {
		return nil, err
	}

	// 6. 缓存任务日志到Redis
	logKey := "task_log:" + string(rune(userTaskLog.Id))
	l.svcCtx.Redis.Setex(logKey, 3600, userTaskLog) // 缓存1小时

	// 7. 发布任务日志创建事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.log.created", map[string]interface{}{
		"log_id":        userTaskLog.Id,
		"user_id":       userTaskLog.UserId,
		"task_id":       userTaskLog.TaskId,
		"progress":      userTaskLog.Progress,
		"created_at":    now,
	})

	// 8. 如果进度为100%，发送完成通知邮件
	if req.Progress == 100 {
		for _, userId := range task.ResponsibleUserIds {
			if userId != req.User_id {
				user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
				l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "任务完成通知", 
					fmt.Sprintf("任务 %s 已完成！", task.Title))
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
