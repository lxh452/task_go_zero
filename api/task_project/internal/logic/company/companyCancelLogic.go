// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package company

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompanyCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注销企业
func NewCompanyCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompanyCancelLogic {
	return &CompanyCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompanyCancelLogic) CompanyCancel(req *types.IdReq) (resp *types.ActionResp, err error) {
	// 1. 查找公司
	company, err := l.svcCtx.CompanyModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("公司不存在")
		}
		return nil, err
	}

	// 2. 检查公司状态
	if company.Status != 1 {
		return nil, logx.Errorf("公司已注销")
	}

	// 3. 软删除公司（设置状态为0）
	now := time.Now()
	company.Status = 0
	company.UpdatedAt = now

	err = l.svcCtx.CompanyModel.Update(l.ctx, company)
	if err != nil {
		return nil, err
	}

	// 4. 删除Redis缓存
	companyKey := "company:" + string(rune(company.Id))
	l.svcCtx.Redis.Del(companyKey)

	// 5. 发布公司注销事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "company.cancelled", map[string]interface{}{
		"company_id":   company.Id,
		"name":         company.Name,
		"cancelled_at": now,
	})

	// 6. 发送注销通知邮件给所有员工
	// 这里可以查询该公司的所有员工并发送邮件
	l.svcCtx.Mailer.SendEmail(l.ctx, company.Email, "公司注销通知", 
		fmt.Sprintf("公司 %s 已注销，请联系管理员。", company.Name))

	return &types.ActionResp{Success: true, Message: "公司注销成功"}, nil
}
