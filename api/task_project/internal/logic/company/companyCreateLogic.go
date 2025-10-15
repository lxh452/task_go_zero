// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package company

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建企业
func NewCompanyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyCreateLogic {
	return &CompanyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyCreateLogic) CompanyCreate(req *types.CompanyCreateReq) (resp *types.Company, err error) {
	// 1. 检查公司名称是否已存在
	_, err = l.svcCtx.CompanyModel.FindOneByName(l.ctx, req.Name)
	if err == nil {
		return nil, logx.Errorf("公司名称已存在: %s", req.Name)
	}
	if err != core.ErrNotFound {
		return nil, err
	}

	// 2. 创建公司
	now := time.Now()
	company := &core.Company{
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		Status:      1, // 正常
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = l.svcCtx.CompanyModel.Insert(l.ctx, company)
	if err != nil {
		return nil, err
	}

	// 3. 缓存公司信息到Redis
	companyKey := "company:" + string(rune(company.Id))
	l.svcCtx.Redis.Setex(companyKey, 3600, company) // 缓存1小时

	// 4. 发布公司创建事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "company.created", map[string]interface{}{
		"company_id":   company.Id,
		"name":         company.Name,
		"created_at":   now,
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
