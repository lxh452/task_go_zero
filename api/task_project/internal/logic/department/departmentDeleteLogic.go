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
	// 1. 查找部门
	department, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("部门不存在")
		}
		return nil, err
	}

	// 2. 检查部门状态
	if department.Status != 1 {
		return nil, logx.Errorf("部门已删除")
	}

	// 3. 检查是否有子部门
	// 这里可以添加检查子部门的逻辑，如果有子部门则不允许删除

	// 4. 软删除部门（设置状态为0）
	now := time.Now()
	department.Status = 0
	department.UpdatedAt = now

	err = l.svcCtx.DepartmentModel.Update(l.ctx, department)
	if err != nil {
		return nil, err
	}

	// 5. 删除Redis缓存
	deptKey := "department:" + string(rune(department.Id))
	l.svcCtx.Redis.Del(deptKey)

	// 6. 发布部门删除事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "department.deleted", map[string]interface{}{
		"department_id": department.Id,
		"company_id":    department.CompanyId,
		"name":          department.Name,
		"deleted_at":    now,
	})

	// 7. 发送删除通知邮件
	// 这里可以查询该部门的所有员工并发送邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, "", "部门删除通知", 
		fmt.Sprintf("部门 %s 已删除，请联系管理员。", department.Name))

	return &types.ActionResp{Success: true, Message: "部门删除成功"}, nil
}
