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
	// 1. 查找部门
	department, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("部门不存在")
		}
		return nil, err
	}

	// 2. 验证父部门是否存在（如果父部门有变化）
	if req.Parent_id != 0 && req.Parent_id != department.ParentId {
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
		if parentDept.CompanyId != department.CompanyId {
			return nil, logx.Errorf("父部门不属于同一公司")
		}
		// 检查是否会造成循环引用
		if req.Parent_id == req.Id {
			return nil, logx.Errorf("不能将自己设为父部门")
		}
	}

	// 3. 更新部门信息
	now := time.Now()
	if req.Name != "" {
		department.Name = req.Name
	}
	if req.Description != "" {
		department.Description = req.Description
	}
	if req.Parent_id != 0 {
		department.ParentId = req.Parent_id
	}
	department.UpdatedAt = now

	err = l.svcCtx.DepartmentModel.Update(l.ctx, department)
	if err != nil {
		return nil, err
	}

	// 4. 更新Redis缓存
	deptKey := "department:" + string(rune(department.Id))
	l.svcCtx.Redis.Setex(deptKey, 3600, department) // 缓存1小时

	// 5. 发布部门更新事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "department.updated", map[string]interface{}{
		"department_id": department.Id,
		"company_id":    department.CompanyId,
		"parent_id":     department.ParentId,
		"name":          department.Name,
		"updated_at":    now,
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
