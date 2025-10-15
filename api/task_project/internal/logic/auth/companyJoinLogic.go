// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyJoinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请加入企业
func NewCompanyJoinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyJoinLogic {
	return &CompanyJoinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyJoinLogic) CompanyJoin(req *types.CompanyJoinReq) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
