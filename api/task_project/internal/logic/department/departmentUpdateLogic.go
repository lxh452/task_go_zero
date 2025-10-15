// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package department

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改部门
func NewDepartmentUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentUpdateLogic {
	return &DepartmentUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentUpdateLogic) DepartmentUpdate(req *types.DepartmentUpdateReq) (resp *types.Department, err error) {
	// todo: add your logic here and delete this line

	return
}
