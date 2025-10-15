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
	// todo: add your logic here and delete this line

	return
}
