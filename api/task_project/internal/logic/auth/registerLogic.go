// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册并创建账号
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.UserAccount, err error) {
	// 1. 检查账号是否已存在
	_, err = l.svcCtx.UserAccountModel.FindOneByAccount(l.ctx, req.Account)
	if err == nil {
		return nil, logx.Errorf("账号已存在: %s", req.Account)
	}
	if err != core.ErrNotFound {
		return nil, err
	}

	// 2. 验证公司和部门是否存在
	company, err := l.svcCtx.CompanyModel.FindOne(l.ctx, req.Company_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("公司不存在")
		}
		return nil, err
	}
	if company.Status != 1 {
		return nil, logx.Errorf("公司已禁用")
	}

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

	// 3. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户账号
	now := time.Now()
	userAccount := &core.UserAccount{
		CompanyId:    req.Company_id,
		DepartmentId: req.Department_id,
		Account:      req.Account,
		Name:         req.Name,
		Email:        req.Email,
		RoleTags:     req.Role_tags,
		Status:       1, // 在职
		HiredAt:      &now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. 创建认证账号
	authAccount := &core.AuthAccount{
		UserId:          0, // 将在插入后设置
		PasswordHash:    string(hashedPassword),
		LoginFailedCount: 0,
	}

	// 6. 事务插入
	err = l.svcCtx.UserAccountModel.Insert(l.ctx, userAccount)
	if err != nil {
		return nil, err
	}

	authAccount.UserId = userAccount.Id
	err = l.svcCtx.AuthAccountModel.Insert(l.ctx, authAccount)
	if err != nil {
		return nil, err
	}

	// 7. 缓存用户信息到Redis
	userKey := "user:" + string(rune(userAccount.Id))
	l.svcCtx.Redis.Setex(userKey, int(l.svcCtx.Config.Auth.AccessExpire), userAccount)

	// 8. 发布注册事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "user.register", map[string]interface{}{
		"user_id":      userAccount.Id,
		"account":      userAccount.Account,
		"company_id":   userAccount.CompanyId,
		"department_id": userAccount.DepartmentId,
		"registered_at": now,
	})

	// 9. 发送欢迎邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, userAccount.Email, "欢迎加入企业任务管理系统", 
		fmt.Sprintf("欢迎 %s 加入 %s！", userAccount.Name, company.Name))

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
