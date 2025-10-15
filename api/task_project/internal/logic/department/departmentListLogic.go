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
	// 1. 查询部门列表
	departments, err := l.svcCtx.DepartmentModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为响应格式
	var departmentList []types.Department
	for _, department := range departments {
		departmentList = append(departmentList, types.Department{
			Id:           department.Id,
			Company_id:   department.CompanyId,
			Parent_id:    department.ParentId,
			Name:         department.Name,
			Description:  department.Description,
			Status:       department.Status,
			Created_at:   department.CreatedAt.String(),
			Updated_at:   department.UpdatedAt.String(),
		})
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "department.list_queried", map[string]interface{}{
		"page":      req.Page,
		"size":      req.Size,
		"count":     len(departmentList),
		"queried_at": time.Now(),
	})

	return &types.DepartmentListResp{
		List:  departmentList,
		Total: int64(len(departmentList)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
