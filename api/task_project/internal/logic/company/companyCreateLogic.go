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
	// todo: add your logic here and delete this line

	return
}
