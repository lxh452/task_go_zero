// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建任务
func NewTaskCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCreateLogic {
	return &TaskCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskCreateLogic) TaskCreate(req *types.TaskCreateReq) (resp *types.Task, err error) {
	// 1. 验证公司和部门是否存在
	company, err := l.svcCtx.CompanyModel.FindOne(l.ctx, req.Company_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("公司不存在")
		}
		return nil, err
	}
	if company.Status != 1 {
		return nil, logx.Errorf("公司已禁用")
	}

	department, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, req.Department_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("部门不存在")
		}
		return nil, err
	}
	if department.Status != 1 {
		return nil, logx.Errorf("部门已禁用")
	}

	// 2. 验证负责人是否存在
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
		if user.CompanyId != req.Company_id {
			return nil, logx.Errorf("负责人不属于指定公司: %d", userId)
		}
	}

	// 3. 验证节点用户是否存在
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
		if user.CompanyId != req.Company_id {
			return nil, logx.Errorf("节点用户不属于指定公司: %d", userId)
		}
	}

	// 4. 创建任务
	now := time.Now()
	task := &core.Task{
		CompanyId:           req.Company_id,
		DepartmentId:        req.Department_id,
		ResponsibleUserIds:  req.Responsible_user_ids,
		NodeUserIds:         req.Node_user_ids,
		NodeKey:             req.Node_key,
		Title:               req.Title,
		Description:         req.Description,
		AttachmentUrl:       req.Attachment_url,
		Status:              1, // 进行中
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	err = l.svcCtx.TaskModel.Insert(l.ctx, task)
	if err != nil {
		return nil, err
	}

	// 5. 缓存任务信息到Redis
	taskKey := "task:" + string(rune(task.Id))
	l.svcCtx.Redis.Setex(taskKey, 3600, task) // 缓存1小时

	// 6. 发布任务创建事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.created", map[string]interface{}{
		"task_id":            task.Id,
		"company_id":         task.CompanyId,
		"department_id":      task.DepartmentId,
		"responsible_users":  task.ResponsibleUserIds,
		"node_users":         task.NodeUserIds,
		"created_at":         now,
	})

	// 7. 发送通知邮件给负责人和节点用户
	for _, userId := range req.Responsible_user_ids {
		user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
		l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "新任务分配", 
			fmt.Sprintf("您有一个新任务：%s", task.Title))
	}

	for _, userId := range req.Node_user_ids {
		user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
		l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "新任务节点", 
			fmt.Sprintf("您有一个新任务节点：%s", task.Title))
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
