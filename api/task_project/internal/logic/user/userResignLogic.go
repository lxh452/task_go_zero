// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserResignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户离职
func NewUserResignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserResignLogic {
	return &UserResignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserResignLogic) UserResign(req *types.UserResignReq) (resp *types.ActionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
