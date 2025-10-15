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
	// todo: add your logic here and delete this line

	return
}
