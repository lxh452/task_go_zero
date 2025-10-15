// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package company

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注销企业
func NewCompanyCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyCancelLogic {
	return &CompanyCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyCancelLogic) CompanyCancel(req *types.IdReq) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
