// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandoverListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 交接记录列表
func NewHandoverListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandoverListLogic {
	return &HandoverListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandoverListLogic) HandoverList(req *types.PageReq) (resp *types.HandoverListResp, err error) {
	// 1. 查询交接记录列表（使用UserTaskLog模拟）
	logs, err := l.svcCtx.UserTaskLogModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 过滤交接记录（handover_id不为0的记录）
	var handoverList []types.HandoverRecord
	for _, log := range logs {
		if log.HandoverId != 0 {
			handoverList = append(handoverList, types.HandoverRecord{
				Id:              log.Id,
				From_user_id:    log.CollaboratorId,
				To_user_id:      log.UserId,
				Task_id:         log.TaskId,
				Node_id:         log.NodeId,
				Handover_id:     log.HandoverId,
				Status:          l.getHandoverStatus(log.Progress),
				Created_at:      log.CreatedAt.String(),
				Updated_at:      log.UpdatedAt.String(),
			})
		}
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "handover.list_queried", map[string]interface{}{
		"page":      req.Page,
		"size":      req.Size,
		"count":     len(handoverList),
		"queried_at": time.Now(),
	})

	return &types.HandoverListResp{
		List:  handoverList,
		Total: int64(len(handoverList)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// 获取交接状态
func (l *HandoverListLogic) getHandoverStatus(progress int) string {
	switch progress {
	case -1:
		return "rejected"
	case 0:
		return "pending"
	case 100:
		return "completed"
	default:
		return "unknown"
	}
}
