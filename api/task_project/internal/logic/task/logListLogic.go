// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 任务明细列表
func NewLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogListLogic {
	return &LogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogListLogic) LogList(req *types.PageReq) (resp *types.LogListResp, err error) {
	// 1. 查询任务日志列表
	logs, err := l.svcCtx.UserTaskLogModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为响应格式
	var logList []types.UserTaskLog
	for _, log := range logs {
		logList = append(logList, types.UserTaskLog{
			Id:              log.Id,
			User_id:         log.UserId,
			Task_id:         log.TaskId,
			Node_id:         log.NodeId,
			Collaborator_id: log.CollaboratorId,
			Progress:        log.Progress,
			Due_date:        log.DueDate.String(),
			Handover_id:     log.HandoverId,
			Created_at:      log.CreatedAt.String(),
			Updated_at:      log.UpdatedAt.String(),
		})
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.log.list_queried", map[string]interface{}{
		"page":      req.Page,
		"size":      req.Size,
		"count":     len(logList),
		"queried_at": time.Now(),
	})

	return &types.LogListResp{
		List:  logList,
		Total: int64(len(logList)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
