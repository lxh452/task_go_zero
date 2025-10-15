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
	// todo: add your logic here and delete this line

	return
}
