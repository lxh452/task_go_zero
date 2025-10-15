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
	// todo: add your logic here and delete this line

	return
}
