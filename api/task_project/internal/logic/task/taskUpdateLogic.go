// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新任务
func NewTaskUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskUpdateLogic {
	return &TaskUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskUpdateLogic) TaskUpdate(req *types.TaskUpdateReq) (resp *types.Task, err error) {
	// 1. 查找任务
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务不存在")
		}
		return nil, err
	}

	// 2. 验证负责人是否存在（如果有变化）
	if req.Responsible_user_ids != nil {
		for _, userId := range req.Responsible_user_ids {
			user, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
			if err != nil {
				if err == core.ErrNotFound {
					return nil, logx.Errorf("负责人不存在: %d", userId)
				}
				return nil, err
			}
			if user.Status != 1 {
				return nil, logx.Errorf("负责人不在职: %d", userId)
			}
			if user.CompanyId != task.CompanyId {
				return nil, logx.Errorf("负责人不属于同一公司: %d", userId)
			}
		}
	}

	// 3. 验证节点用户是否存在（如果有变化）
	if req.Node_user_ids != nil {
		for _, userId := range req.Node_user_ids {
			user, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
			if err != nil {
				if err == core.ErrNotFound {
					return nil, logx.Errorf("节点用户不存在: %d", userId)
				}
				return nil, err
			}
			if user.Status != 1 {
				return nil, logx.Errorf("节点用户不在职: %d", userId)
			}
			if user.CompanyId != task.CompanyId {
				return nil, logx.Errorf("节点用户不属于同一公司: %d", userId)
			}
		}
	}

	// 4. 更新任务信息
	now := time.Now()
	oldStatus := task.Status
	
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Attachment_url != "" {
		task.AttachmentUrl = req.Attachment_url
	}
	if req.Responsible_user_ids != nil {
		task.ResponsibleUserIds = req.Responsible_user_ids
	}
	if req.Node_user_ids != nil {
		task.NodeUserIds = req.Node_user_ids
	}
	if req.Node_key != "" {
		task.NodeKey = req.Node_key
	}
	if req.Status != 0 {
		task.Status = req.Status
	}
	task.UpdatedAt = now

	err = l.svcCtx.TaskModel.Update(l.ctx, task)
	if err != nil {
		return nil, err
	}

	// 5. 更新Redis缓存
	taskKey := "task:" + string(rune(task.Id))
	l.svcCtx.Redis.Setex(taskKey, 3600, task) // 缓存1小时

	// 6. 发布任务更新事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.updated", map[string]interface{}{
		"task_id":            task.Id,
		"company_id":         task.CompanyId,
		"department_id":      task.DepartmentId,
		"responsible_users":  task.ResponsibleUserIds,
		"node_users":         task.NodeUserIds,
		"updated_at":         now,
	})

	// 7. 如果状态有变化，发送通知邮件
	if oldStatus != task.Status {
		for _, userId := range task.ResponsibleUserIds {
			user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
			l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "任务状态更新", 
				fmt.Sprintf("任务 %s 状态已更新为 %d", task.Title, task.Status))
		}
	}

	return &types.Task{
		Id:                   task.Id,
		Company_id:           task.CompanyId,
		Department_id:        task.DepartmentId,
		Responsible_user_ids: task.ResponsibleUserIds,
		Node_user_ids:        task.NodeUserIds,
		Node_key:             task.NodeKey,
		Title:                task.Title,
		Description:          task.Description,
		Attachment_url:       task.AttachmentUrl,
		Status:               task.Status,
		Created_at:           task.CreatedAt.String(),
		Updated_at:           task.UpdatedAt.String(),
	}, nil
}
