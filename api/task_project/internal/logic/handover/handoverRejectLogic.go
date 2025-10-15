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
	// todo: add your logic here and delete this line

	return
}
