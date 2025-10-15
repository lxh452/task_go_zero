// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandoverProposeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发起交接提议
func NewHandoverProposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandoverProposeLogic {
	return &HandoverProposeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandoverProposeLogic) HandoverPropose(req *types.HandoverProposeReq) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
