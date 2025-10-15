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

	// 4. 创建用户账号（待审核状态）
	now := time.Now()
	userAccount := &core.UserAccount{
		CompanyId:    req.Company_id,
		DepartmentId: req.Department_id,
		Account:      req.Account,
		Name:         req.Name,
		Email:        req.Email,
		RoleTags:     req.Role_tags,
		Status:       0, // 待审核
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

	// 7. 发布加入企业事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "user.join_company", map[string]interface{}{
		"user_id":      userAccount.Id,
		"account":      userAccount.Account,
		"company_id":   userAccount.CompanyId,
		"department_id": userAccount.DepartmentId,
		"joined_at":    now,
	})

	// 8. 发送申请邮件给公司管理员
	l.svcCtx.Mailer.SendEmail(l.ctx, company.Email, "新员工加入申请", 
		fmt.Sprintf("用户 %s 申请加入公司 %s，请审核。", userAccount.Name, company.Name))

	return &types.ActionResp{Success: true, Message: "申请已提交，等待审核"}, nil
}
