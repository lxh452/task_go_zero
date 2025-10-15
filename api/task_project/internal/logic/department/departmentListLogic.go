// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package department

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 部门列表
func NewDepartmentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentListLogic {
	return &DepartmentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentListLogic) DepartmentList(req *types.PageReq) (resp *types.DepartmentListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
