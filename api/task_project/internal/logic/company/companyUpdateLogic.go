// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package company

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改企业信息/状态
func NewCompanyUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyUpdateLogic {
	return &CompanyUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyUpdateLogic) CompanyUpdate(req *types.CompanyUpdateReq) (resp *types.Company, err error) {
	// 1. 查找公司
	company, err := l.svcCtx.CompanyModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("公司不存在")
		}
		return nil, err
	}

	// 2. 检查公司名称是否冲突（如果名称有变化）
	if req.Name != "" && req.Name != company.Name {
		_, err = l.svcCtx.CompanyModel.FindOneByName(l.ctx, req.Name)
		if err == nil {
			return nil, logx.Errorf("公司名称已存在: %s", req.Name)
		}
		if err != core.ErrNotFound {
			return nil, err
		}
	}

	// 3. 更新公司信息
	now := time.Now()
	if req.Name != "" {
		company.Name = req.Name
	}
	if req.Description != "" {
		company.Description = req.Description
	}
	if req.Email != "" {
		company.Email = req.Email
	}
	if req.Phone != "" {
		company.Phone = req.Phone
	}
	if req.Address != "" {
		company.Address = req.Address
	}
	company.UpdatedAt = now

	err = l.svcCtx.CompanyModel.Update(l.ctx, company)
	if err != nil {
		return nil, err
	}

	// 4. 更新Redis缓存
	companyKey := "company:" + string(rune(company.Id))
	l.svcCtx.Redis.Setex(companyKey, 3600, company) // 缓存1小时

	// 5. 发布公司更新事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "company.updated", map[string]interface{}{
		"company_id":   company.Id,
		"name":         company.Name,
		"updated_at":   now,
	})

	return &types.Company{
		Id:          company.Id,
		Name:        company.Name,
		Description: company.Description,
		Email:       company.Email,
		Phone:       company.Phone,
		Address:     company.Address,
		Status:      company.Status,
		Created_at:  company.CreatedAt.String(),
		Updated_at:  company.UpdatedAt.String(),
	}, nil
}
