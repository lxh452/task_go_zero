// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package department

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除部门
func NewDepartmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentDeleteLogic {
	return &DepartmentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentDeleteLogic) DepartmentDelete(req *types.IdReq) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
