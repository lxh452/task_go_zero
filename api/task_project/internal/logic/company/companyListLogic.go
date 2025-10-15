// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package company

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 企业列表
func NewCompanyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyListLogic {
	return &CompanyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyListLogic) CompanyList(req *types.PageReq) (resp *types.CompanyListResp, err error) {
	// 1. 查询公司列表
	companies, err := l.svcCtx.CompanyModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为响应格式
	var companyList []types.Company
	for _, company := range companies {
		companyList = append(companyList, types.Company{
			Id:          company.Id,
			Name:        company.Name,
			Description: company.Description,
			Email:       company.Email,
			Phone:       company.Phone,
			Address:     company.Address,
			Status:      company.Status,
			Created_at:  company.CreatedAt.String(),
			Updated_at:  company.UpdatedAt.String(),
		})
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "company.list_queried", map[string]interface{}{
		"page":      req.Page,
		"size":      req.Size,
		"count":     len(companyList),
		"queried_at": time.Now(),
	})

	return &types.CompanyListResp{
		List:  companyList,
		Total: int64(len(companyList)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
