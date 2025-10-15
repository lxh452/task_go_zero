// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dispatch

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 确认派发结果
func NewDispatchConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchConfirmLogic {
	return &DispatchConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchConfirmLogic) DispatchConfirm(req *types.DispatchConfirmReq) (resp *types.ActionResp, err error) {
	// 1. 查找任务
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Task_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务不存在")
		}
		return nil, err
	}

	// 2. 验证选中的用户是否存在且在职
	for _, userId := range req.Selected_user_ids {
		user, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
		if err != nil {
			if err == core.ErrNotFound {
				return nil, logx.Errorf("用户不存在: %d", userId)
			}
			return nil, err
		}
		if user.Status != 1 {
			return nil, logx.Errorf("用户不在职: %d", userId)
		}
		if user.CompanyId != task.CompanyId {
			return nil, logx.Errorf("用户不属于同一公司: %d", userId)
		}
	}

	// 3. 更新任务负责人
	task.ResponsibleUserIds = req.Selected_user_ids
	task.UpdatedAt = time.Now()

	err = l.svcCtx.TaskModel.Update(l.ctx, task)
	if err != nil {
		return nil, err
	}

	// 4. 更新Redis缓存
	taskKey := "task:" + string(rune(task.Id))
	l.svcCtx.Redis.Setex(taskKey, 3600, task)

	// 5. 发布派发确认事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.dispatch.confirm", map[string]interface{}{
		"task_id":            task.Id,
		"company_id":         task.CompanyId,
		"department_id":      task.DepartmentId,
		"responsible_users":  req.Selected_user_ids,
		"confirmed_at":       time.Now(),
	})

	// 6. 发送通知邮件给选中的用户
	for _, userId := range req.Selected_user_ids {
		user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
		l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "任务分配确认", 
			fmt.Sprintf("您已被确认为任务负责人：%s", task.Title))
	}

	return &types.ActionResp{Success: true, Message: "任务派发确认成功"}, nil
}
