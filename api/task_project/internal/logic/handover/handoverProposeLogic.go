// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandoverProposeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发起交接提议
func NewHandoverProposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandoverProposeLogic {
	return &HandoverProposeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandoverProposeLogic) HandoverPropose(req *types.HandoverProposeReq) (resp *types.ActionResp, err error) {
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

	// 6. 创建交接记录（使用UserTaskLog模拟）
	now := time.Now()
	handoverLog := &core.UserTaskLog{
		UserId:         req.From_user_id,
		TaskId:         req.Task_id,
		NodeId:         req.Node_id,
		CollaboratorId: req.To_user_id,
		Progress:       0, // 交接中
		HandoverId:     req.Handover_id,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = l.svcCtx.UserTaskLogModel.Insert(l.ctx, handoverLog)
	if err != nil {
		return nil, err
	}

	// 7. 缓存交接记录到Redis
	logKey := "handover:" + string(rune(handoverLog.Id))
	l.svcCtx.Redis.Setex(logKey, 3600, handoverLog)

	// 8. 发布交接提议事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.handover.propose", map[string]interface{}{
		"handover_id":   handoverLog.Id,
		"from_user_id":  req.From_user_id,
		"to_user_id":    req.To_user_id,
		"task_id":       req.Task_id,
		"proposed_at":   now,
	})

	// 9. 发送交接提议邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, toUser.Email, "任务交接提议", 
		fmt.Sprintf("用户 %s 提议将任务 %s 交接给您，请确认。", fromUser.Name, task.Title))

	l.svcCtx.Mailer.SendEmail(l.ctx, fromUser.Email, "交接提议已发送", 
		fmt.Sprintf("您已向用户 %s 发送任务 %s 的交接提议。", toUser.Name, task.Title))

	return &types.ActionResp{Success: true, Message: "交接提议已发送"}, nil
}
