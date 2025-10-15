// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dispatch

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 自动派发候选预览
func NewDispatchPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchPreviewLogic {
	return &DispatchPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchPreviewLogic) DispatchPreview(req *types.DispatchPreviewReq) (resp *types.DispatchPreviewResp, err error) {
	// todo: add your logic here and delete this line

	return
}
