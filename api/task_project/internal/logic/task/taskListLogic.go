// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 任务列表
func NewTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskListLogic {
	return &TaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskListLogic) TaskList(req *types.PageReq) (resp *types.TaskListResp, err error) {
	// 1. 查询任务列表
	tasks, err := l.svcCtx.TaskModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为响应格式
	var taskList []types.Task
	for _, task := range tasks {
		taskList = append(taskList, types.Task{
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
		})
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.list_queried", map[string]interface{}{
		"page":      req.Page,
		"size":      req.Size,
		"count":     len(taskList),
		"queried_at": time.Now(),
	})

	return &types.TaskListResp{
		List:  taskList,
		Total: int64(len(taskList)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
