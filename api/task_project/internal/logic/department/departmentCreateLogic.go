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
	// 1. 验证公司是否存在
	company, err := l.svcCtx.CompanyModel.FindOne(l.ctx, req.Company_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("公司不存在")
		}
		return nil, err
	}
	if company.Status != 1 {
		return nil, logx.Errorf("公司已禁用")
	}

	// 2. 验证父部门是否存在（如果有父部门）
	if req.Parent_id != 0 {
		parentDept, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, req.Parent_id)
		if err != nil {
			if err == core.ErrNotFound {
				return nil, logx.Errorf("父部门不存在")
			}
			return nil, err
		}
		if parentDept.Status != 1 {
			return nil, logx.Errorf("父部门已禁用")
		}
		if parentDept.CompanyId != req.Company_id {
			return nil, logx.Errorf("父部门不属于指定公司")
		}
	}

	// 3. 创建部门
	now := time.Now()
	department := &core.Department{
		CompanyId:  req.Company_id,
		ParentId:   req.Parent_id,
		Name:       req.Name,
		Description: req.Description,
		Status:     1, // 正常
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = l.svcCtx.DepartmentModel.Insert(l.ctx, department)
	if err != nil {
		return nil, err
	}

	// 4. 缓存部门信息到Redis
	deptKey := "department:" + string(rune(department.Id))
	l.svcCtx.Redis.Setex(deptKey, 3600, department) // 缓存1小时

	// 5. 发布部门创建事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "department.created", map[string]interface{}{
		"department_id": department.Id,
		"company_id":    department.CompanyId,
		"parent_id":     department.ParentId,
		"name":          department.Name,
		"created_at":    now,
	})

	return &types.Department{
		Id:           department.Id,
		Company_id:   department.CompanyId,
		Parent_id:    department.ParentId,
		Name:         department.Name,
		Description:  department.Description,
		Status:       department.Status,
		Created_at:   department.CreatedAt.String(),
		Updated_at:   department.UpdatedAt.String(),
	}, nil
}
