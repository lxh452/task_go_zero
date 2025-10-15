// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"time"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"
	"task_Project/model/core"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 查找用户账号
	userAccount, err := l.svcCtx.UserAccountModel.FindOneByAccount(l.ctx, req.Account)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("用户不存在: %s", req.Account)
		}
		return nil, err
	}

	// 2. 检查用户状态
	if userAccount.Status != 1 {
		return nil, logx.Errorf("用户账号已禁用")
	}

	// 3. 验证密码
	authAccount, err := l.svcCtx.AuthAccountModel.FindOne(l.ctx, userAccount.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("认证信息不存在")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(authAccount.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, logx.Errorf("密码错误")
	}

	// 4. 生成JWT token
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(now),
		Subject:   string(rune(userAccount.Id)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
	if err != nil {
		return nil, err
	}

	// 5. 更新最后登录时间
	authAccount.LastLoginAt = &now
	l.svcCtx.AuthAccountModel.Update(l.ctx, authAccount)

	// 6. 缓存用户信息到Redis
	userKey := "user:" + string(rune(userAccount.Id))
	l.svcCtx.Redis.Setex(userKey, int(l.svcCtx.Config.Auth.AccessExpire), userAccount)

	// 7. 发布登录事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "user.login", map[string]interface{}{
		"user_id":  userAccount.Id,
		"account":  userAccount.Account,
		"login_at": now,
	})

	return &types.LoginResp{
		Token: signed,
		User: types.UserAccount{
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
		},
	}, nil
}
