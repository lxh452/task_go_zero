// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package schedule

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScheduleQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询任务排期
func NewScheduleQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleQueryLogic {
	return &ScheduleQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleQueryLogic) ScheduleQuery(req *types.ScheduleQueryReq) (resp *types.LogListResp, err error) {
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

	// 2. 查询用户的任务日志
	logs, err := l.svcCtx.UserTaskLogModel.FindAll(l.ctx, 1, 1000)
	if err != nil {
		return nil, err
	}

	// 3. 过滤用户相关的任务日志
	var userLogs []types.UserTaskLog
	for _, log := range logs {
		if log.UserId == req.User_id {
			// 根据时间范围过滤
			if l.isLogInTimeRange(log, req.Start_time, req.End_time, req.Granularity) {
				userLogs = append(userLogs, types.UserTaskLog{
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
		}
	}

	// 4. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "schedule.queried", map[string]interface{}{
		"user_id":      req.User_id,
		"start_time":   req.Start_time,
		"end_time":     req.End_time,
		"granularity":  req.Granularity,
		"log_count":    len(userLogs),
		"queried_at":   time.Now(),
	})

	return &types.LogListResp{
		List:  userLogs,
		Total: int64(len(userLogs)),
		Page:  1,
		Size:  int64(len(userLogs)),
	}, nil
}

// 检查任务日志是否在时间范围内
func (l *ScheduleQueryLogic) isLogInTimeRange(log *core.UserTaskLog, startTime, endTime time.Time, granularity string) bool {
	logTime := log.CreatedAt
	
	switch granularity {
	case "day":
		return logTime.After(startTime) && logTime.Before(endTime)
	case "month":
		// 检查是否在同一个月
		return logTime.Year() == startTime.Year() && logTime.Month() == startTime.Month()
	case "year":
		// 检查是否在同一年
		return logTime.Year() == startTime.Year()
	default:
		return true
	}
}
