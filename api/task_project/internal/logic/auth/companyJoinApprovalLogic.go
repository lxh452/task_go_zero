// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyJoinApprovalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 审批加入企业申请
func NewCompanyJoinApprovalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyJoinApprovalLogic {
	return &CompanyJoinApprovalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyJoinApprovalLogic) CompanyJoinApproval(req *types.CompanyJoinApprovalReq) (resp *types.ActionResp, err error) {
	// 1. 查找用户账号
	userAccount, err := l.svcCtx.UserAccountModel.FindOne(l.ctx, req.User_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在")
		}
		return nil, err
	}

	// 2. 检查用户状态（必须是待审核状态）
	if userAccount.Status != 0 {
		return nil, logx.Errorf("用户状态不是待审核")
	}

	// 3. 更新用户状态
	if req.Approved {
		userAccount.Status = 1 // 在职
	} else {
		userAccount.Status = -1 // 拒绝
	}
	userAccount.UpdatedAt = time.Now()

	err = l.svcCtx.UserAccountModel.Update(l.ctx, userAccount)
	if err != nil {
		return nil, err
	}

	// 4. 缓存用户信息到Redis
	userKey := "user:" + string(rune(userAccount.Id))
	if req.Approved {
		l.svcCtx.Redis.Setex(userKey, int(l.svcCtx.Config.Auth.AccessExpire), userAccount)
	} else {
		l.svcCtx.Redis.Del(userKey)
	}

	// 5. 发布审核结果事件到MQ
	eventType := "user.join_rejected"
	if req.Approved {
		eventType = "user.join_approved"
	}
	l.svcCtx.MQ.Publish(l.ctx, eventType, map[string]interface{}{
		"user_id":    userAccount.Id,
		"account":    userAccount.Account,
		"approved":   req.Approved,
		"approved_at": time.Now(),
	})

	// 6. 发送审核结果邮件
	if req.Approved {
		l.svcCtx.Mailer.SendEmail(l.ctx, userAccount.Email, "加入企业申请已通过", 
			fmt.Sprintf("恭喜 %s，您的加入企业申请已通过！", userAccount.Name))
	} else {
		l.svcCtx.Mailer.SendEmail(l.ctx, userAccount.Email, "加入企业申请被拒绝", 
			fmt.Sprintf("很抱歉 %s，您的加入企业申请被拒绝。", userAccount.Name))
	}

	message := "申请已拒绝"
	if req.Approved {
		message = "申请已通过"
	}

	return &types.ActionResp{Success: true, Message: message}, nil
}
