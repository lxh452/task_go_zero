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
	// 1. 查找用户
	userAccount, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.User_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在")
		}
		return nil, err
	}

	// 2. 检查用户状态
	if userAccount.Status != 1 {
		return nil, logx.Errorf("用户不在职")
	}

	// 3. 设置用户为离职状态
	now := time.Now()
	userAccount.Status = 0 // 离职
	userAccount.LeftAt = &now
	userAccount.UpdatedAt = now

	err = l.svcCtx.UserAccountModel.Update(l.ctx, userAccount)
	if err != nil {
		return nil, err
	}

	// 4. 更新Redis缓存
	userKey := "user:" + string(rune(userAccount.Id))
	l.svcCtx.Redis.Setex(userKey, int(l.svcCtx.Config.Auth.AccessExpire), userAccount)

	// 5. 发布用户离职事件到MQ（触发任务交接）
	l.svcCtx.MQ.Publish(l.ctx, "user.resigned", map[string]interface{}{
		"user_id":      userAccount.Id,
		"account":      userAccount.Account,
		"company_id":   userAccount.CompanyId,
		"department_id": userAccount.DepartmentId,
		"resigned_at":  now,
	})

	// 6. 发送离职通知邮件给用户和部门经理
	l.svcCtx.Mailer.SendEmail(l.ctx, userAccount.Email, "离职通知", 
		fmt.Sprintf("您已从 %s 离职，感谢您的工作！", userAccount.Name))

	// 这里可以查询部门经理并发送通知邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, "", "员工离职通知", 
		fmt.Sprintf("员工 %s 已离职，请处理相关任务交接。", userAccount.Name))

	return &types.ActionResp{Success: true, Message: "用户离职成功"}, nil
}
