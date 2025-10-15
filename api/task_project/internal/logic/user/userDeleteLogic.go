// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除用户
func NewUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeleteLogic {
	return &UserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeleteLogic) UserDelete(req *types.IdReq) (resp *types.ActionResp, err error) {
	// 1. 查找用户
	userAccount, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在")
		}
		return nil, err
	}

	// 2. 检查用户状态
	if userAccount.Status == -1 {
		return nil, logx.Errorf("用户已删除")
	}

	// 3. 软删除用户（设置状态为-1）
	now := time.Now()
	userAccount.Status = -1
	userAccount.LeftAt = &now
	userAccount.UpdatedAt = now

	err = l.svcCtx.UserAccountModel.Update(l.ctx, userAccount)
	if err != nil {
		return nil, err
	}

	// 4. 删除Redis缓存
	userKey := "user:" + string(rune(userAccount.Id))
	l.svcCtx.Redis.Del(userKey)

	// 5. 发布用户删除事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "user.deleted", map[string]interface{}{
		"user_id":      userAccount.Id,
		"account":      userAccount.Account,
		"company_id":   userAccount.CompanyId,
		"department_id": userAccount.DepartmentId,
		"deleted_at":   now,
	})

	// 6. 发送删除通知邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, userAccount.Email, "账号删除通知", 
		fmt.Sprintf("您的账号 %s 已被删除，请联系管理员。", userAccount.Account))

	return &types.ActionResp{Success: true, Message: "用户删除成功"}, nil
}
