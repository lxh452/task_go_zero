// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户档案列表
func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.PageReq) (resp *types.UserListResp, err error) {
	// 1. 查询用户列表
	users, err := l.svcCtx.UserAccountModel.FindAll(l.ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 2. 转换为响应格式
	var userList []types.UserAccount
	for _, user := range users {
		userList = append(userList, types.UserAccount{
			Id:            user.Id,
			Company_id:    user.CompanyId,
			Department_id: user.DepartmentId,
			Account:       user.Account,
			Name:          user.Name,
			Email:         user.Email,
			Role_tags:     user.RoleTags,
			Status:        user.Status,
			Hired_at:      user.HiredAt.String(),
			Left_at:       user.LeftAt.String(),
			Created_at:    user.CreatedAt.String(),
			Updated_at:    user.UpdatedAt.String(),
		})
	}

	// 3. 发布查询事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "user.list_queried", map[string]interface{}{
		"page":      req.Page,
		"size":      req.Size,
		"count":     len(userList),
		"queried_at": time.Now(),
	})

	return &types.UserListResp{
		List:  userList,
		Total: int64(len(userList)),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
