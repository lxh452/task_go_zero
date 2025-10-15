// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandoverRejectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 拒绝交接
func NewHandoverRejectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandoverRejectLogic {
	return &HandoverRejectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandoverRejectLogic) HandoverReject(req *types.HandoverActionReq) (resp *types.ActionResp, err error) {
	// 1. 验证接收用户是否存在且在职
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

	// 2. 验证任务是否存在
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Task_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务不存在")
		}
		return nil, err
	}

	// 3. 创建交接拒绝记录
	now := time.Now()
	handoverLog := &core.UserTaskLog{
		UserId:         req.To_user_id,
		TaskId:         req.Task_id,
		NodeId:         req.Node_id,
		CollaboratorId: req.From_user_id,
		Progress:       -1, // 交接拒绝
		HandoverId:     req.Handover_id,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = l.svcCtx.UserTaskLogModel.Insert(l.ctx, handoverLog)
	if err != nil {
		return nil, err
	}

	// 4. 缓存交接记录到Redis
	logKey := "handover:" + string(rune(handoverLog.Id))
	l.svcCtx.Redis.Setex(logKey, 3600, handoverLog)

	// 5. 发布交接拒绝事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.handover.reject", map[string]interface{}{
		"handover_id":   handoverLog.Id,
		"from_user_id":  req.From_user_id,
		"to_user_id":    req.To_user_id,
		"task_id":       req.Task_id,
		"rejected_at":   now,
	})

	// 6. 发送交接拒绝邮件
	fromUser, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.From_user_id)
	l.svcCtx.Mailer.SendEmail(l.ctx, fromUser.Email, "任务交接被拒绝", 
		fmt.Sprintf("用户 %s 拒绝了任务 %s 的交接提议。", toUser.Name, task.Title))

	l.svcCtx.Mailer.SendEmail(l.ctx, toUser.Email, "交接拒绝已发送", 
		fmt.Sprintf("您已拒绝任务 %s 的交接提议。", task.Title))

	return &types.ActionResp{Success: true, Message: "交接拒绝成功"}, nil
}
