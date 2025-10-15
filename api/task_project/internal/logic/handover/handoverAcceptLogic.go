// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandoverAcceptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 同意交接
func NewHandoverAcceptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandoverAcceptLogic {
	return &HandoverAcceptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandoverAcceptLogic) HandoverAccept(req *types.HandoverActionReq) (resp *types.ActionResp, err error) {
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

	// 3. 更新任务负责人（添加接收用户）
	// 检查接收用户是否已经在负责人列表中
	alreadyResponsible := false
	for _, userId := range task.ResponsibleUserIds {
		if userId == req.To_user_id {
			alreadyResponsible = true
			break
		}
	}

	if !alreadyResponsible {
		task.ResponsibleUserIds = append(task.ResponsibleUserIds, req.To_user_id)
		task.UpdatedAt = time.Now()

		err = l.svcCtx.TaskModel.Update(l.ctx, task)
		if err != nil {
			return nil, err
		}

		// 更新Redis缓存
		taskKey := "task:" + string(rune(task.Id))
		l.svcCtx.Redis.Setex(taskKey, 3600, task)
	}

	// 4. 创建交接完成记录
	now := time.Now()
	handoverLog := &core.UserTaskLog{
		UserId:         req.To_user_id,
		TaskId:         req.Task_id,
		NodeId:         req.Node_id,
		CollaboratorId: req.From_user_id,
		Progress:       100, // 交接完成
		HandoverId:     req.Handover_id,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = l.svcCtx.UserTaskLogModel.Insert(l.ctx, handoverLog)
	if err != nil {
		return nil, err
	}

	// 5. 缓存交接记录到Redis
	logKey := "handover:" + string(rune(handoverLog.Id))
	l.svcCtx.Redis.Setex(logKey, 3600, handoverLog)

	// 6. 发布交接接受事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.handover.accept", map[string]interface{}{
		"handover_id":   handoverLog.Id,
		"from_user_id":  req.From_user_id,
		"to_user_id":    req.To_user_id,
		"task_id":       req.Task_id,
		"accepted_at":   now,
	})

	// 7. 发送交接完成邮件
	fromUser, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.From_user_id)
	l.svcCtx.Mailer.SendEmail(l.ctx, toUser.Email, "任务交接完成", 
		fmt.Sprintf("您已成功接收任务 %s 的交接。", task.Title))

	l.svcCtx.Mailer.SendEmail(l.ctx, fromUser.Email, "任务交接完成", 
		fmt.Sprintf("用户 %s 已接受任务 %s 的交接。", toUser.Name, task.Title))

	return &types.ActionResp{Success: true, Message: "交接接受成功"}, nil
}
