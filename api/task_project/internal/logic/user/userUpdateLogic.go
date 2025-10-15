// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户档案
func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserUpdateReq) (resp *types.UserAccount, err error) {
	// 1. 查找用户
	userAccount, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在")
		}
		return nil, err
	}

	// 2. 检查账号是否冲突（如果账号有变化）
	if req.Account != "" && req.Account != userAccount.Account {
		_, err = l.svcCtx.UserAccountModel.FindOneByAccount(l.ctx, req.Account)
		if err == nil {
			return nil, logx.Errorf("账号已存在: %s", req.Account)
		}
		if err != core.ErrNotFound {
			return nil, err
		}
	}

	// 3. 验证部门是否存在（如果部门有变化）
	if req.Department_id != 0 && req.Department_id != userAccount.DepartmentId {
		department, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, req.Department_id)
		if err != nil {
			if err == core.ErrNotFound {
				return nil, logx.Errorf("部门不存在")
			}
			return nil, err
		}
		if department.Status != 1 {
			return nil, logx.Errorf("部门已禁用")
		}
		if department.CompanyId != userAccount.CompanyId {
			return nil, logx.Errorf("部门不属于同一公司")
		}
	}

	// 4. 更新用户信息
	now := time.Now()
	oldStatus := userAccount.Status
	
	if req.Account != "" {
		userAccount.Account = req.Account
	}
	if req.Name != "" {
		userAccount.Name = req.Name
	}
	if req.Email != "" {
		userAccount.Email = req.Email
	}
	if req.Role_tags != nil {
		userAccount.RoleTags = req.Role_tags
	}
	if req.Department_id != 0 {
		userAccount.DepartmentId = req.Department_id
	}
	if req.Status != 0 {
		userAccount.Status = req.Status
	}
	userAccount.UpdatedAt = now

	err = l.svcCtx.UserAccountModel.Update(l.ctx, userAccount)
	if err != nil {
		return nil, err
	}

	// 5. 更新Redis缓存
	userKey := "user:" + string(rune(userAccount.Id))
	l.svcCtx.Redis.Setex(userKey, int(l.svcCtx.Config.Auth.AccessExpire), userAccount)

	// 6. 发布用户更新事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "user.updated", map[string]interface{}{
		"user_id":      userAccount.Id,
		"account":      userAccount.Account,
		"company_id":   userAccount.CompanyId,
		"department_id": userAccount.DepartmentId,
		"updated_at":   now,
	})

	// 7. 如果状态变为离职，触发任务交接
	if oldStatus == 1 && userAccount.Status == 0 {
		l.svcCtx.MQ.Publish(l.ctx, "user.resigned", map[string]interface{}{
			"user_id":      userAccount.Id,
			"account":      userAccount.Account,
			"company_id":   userAccount.CompanyId,
			"department_id": userAccount.DepartmentId,
			"resigned_at":  now,
		})
	}

	return &types.UserAccount{
		Id:            userAccount.Id,
		Company_id:    userAccount.CompanyId,
		Department_id: userAccount.DepartmentId,
		Account:       userAccount.Account,
		Name:          userAccount.Name,
		Email:         userAccount.Email,
		Role_tags:     userAccount.RoleTags,
		Status:        userAccount.Status,
		Hired_at:      userAccount.HiredAt.String(),
		Left_at:       userAccount.LeftAt.String(),
		Created_at:    userAccount.CreatedAt.String(),
		Updated_at:    userAccount.UpdatedAt.String(),
	}, nil
}
