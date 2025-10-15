// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 直接交接（即时生效）
func NewHandoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandoverLogic {
	return &HandoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandoverLogic) Handover(req *types.HandoverReq) (resp *types.ActionResp, err error) {
	// 1. 验证原用户是否存在且在职
	fromUser, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.From_user_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("原用户不存在")
		}
		return nil, err
	}
	if fromUser.Status != 1 {
		return nil, logx.Errorf("原用户不在职")
	}

	// 2. 验证接收用户是否存在且在职
	toUser, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.To_user_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("接收用户不存在")
		}
		return nil, err
	}
	if toUser.Status != 1 {
		return nil, logx.Errorf("接收用户不在职")
	}

	// 3. 验证任务是否存在
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Task_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务不存在")
		}
		return nil, err
	}

	// 4. 验证原用户是否有权限操作该任务
	hasPermission := false
	for _, userId := range task.ResponsibleUserIds {
		if userId == req.From_user_id {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return nil, logx.Errorf("原用户无权限操作该任务")
	}

	// 5. 验证接收用户是否属于同一公司
	if toUser.CompanyId != fromUser.CompanyId {
		return nil, logx.Errorf("接收用户不属于同一公司")
	}

	// 6. 直接更新任务负责人
	// 移除原用户，添加接收用户
	newResponsibleUsers := make([]int64, 0)
	for _, userId := range task.ResponsibleUserIds {
		if userId != req.From_user_id {
			newResponsibleUsers = append(newResponsibleUsers, userId)
		}
	}
	
	// 检查接收用户是否已经在负责人列表中
	alreadyResponsible := false
	for _, userId := range newResponsibleUsers {
		if userId == req.To_user_id {
			alreadyResponsible = true
			break
		}
	}
	
	if !alreadyResponsible {
		newResponsibleUsers = append(newResponsibleUsers, req.To_user_id)
	}

	task.ResponsibleUserIds = newResponsibleUsers
	task.UpdatedAt = time.Now()

	err = l.svcCtx.TaskModel.Update(l.ctx, task)
	if err != nil {
		return nil, err
	}

	// 7. 更新Redis缓存
	taskKey := "task:" + string(rune(task.Id))
	l.svcCtx.Redis.Setex(taskKey, 3600, task)

	// 8. 创建直接交接记录
	now := time.Now()
	handoverLog := &core.UserTaskLog{
		UserId:         req.To_user_id,
		TaskId:         req.Task_id,
		NodeId:         req.Node_id,
		CollaboratorId: req.From_user_id,
		Progress:       100, // 直接交接完成
		HandoverId:     req.Handover_id,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = l.svcCtx.UserTaskLogModel.Insert(l.ctx, handoverLog)
	if err != nil {
		return nil, err
	}

	// 9. 缓存交接记录到Redis
	logKey := "handover:" + string(rune(handoverLog.Id))
	l.svcCtx.Redis.Setex(logKey, 3600, handoverLog)

	// 10. 发布直接交接事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.handover.direct", map[string]interface{}{
		"handover_id":   handoverLog.Id,
		"from_user_id":  req.From_user_id,
		"to_user_id":    req.To_user_id,
		"task_id":       req.Task_id,
		"handover_at":   now,
	})

	// 11. 发送直接交接邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, toUser.Email, "任务直接交接", 
		fmt.Sprintf("任务 %s 已直接交接给您。", task.Title))

	l.svcCtx.Mailer.SendEmail(l.ctx, fromUser.Email, "任务直接交接", 
		fmt.Sprintf("任务 %s 已直接交接给用户 %s。", task.Title, toUser.Name))

	return &types.ActionResp{Success: true, Message: "直接交接成功"}, nil
}
