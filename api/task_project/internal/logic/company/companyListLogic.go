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
	// todo: add your logic here and delete this line

	return
}
