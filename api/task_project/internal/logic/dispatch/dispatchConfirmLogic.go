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
	// todo: add your logic here and delete this line

	return
}
