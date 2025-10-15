// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package department

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建部门
func NewDepartmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentCreateLogic {
	return &DepartmentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentCreateLogic) DepartmentCreate(req *types.DepartmentCreateReq) (resp *types.Department, err error) {
	// todo: add your logic here and delete this line

	return
}
