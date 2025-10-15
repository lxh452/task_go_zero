// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dispatch

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 自动派发执行
func NewDispatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchLogic {
	return &DispatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchLogic) Dispatch(req *types.DispatchReq) (resp *types.ActionResp, err error) {
	// 1. 查找任务
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Task_id)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, logx.Errorf("任务不存在")
		}
		return nil, err
	}

	// 2. 获取候选用户列表
	candidates, err := l.getCandidates(req.Company_id, req.Department_id, req.Required_role_tags)
	if err != nil {
		return nil, err
	}

	// 3. 计算用户得分并排序
	var scoredUsers []ScoredUser
	for _, user := range candidates {
		score := l.calculateScore(user, req.Required_role_tags, req.Company_id, req.Department_id)
		scoredUsers = append(scoredUsers, ScoredUser{
			User:  user,
			Score: score,
		})
	}

	// 按得分排序（降序）
	sort.Slice(scoredUsers, func(i, j int) bool {
		return scoredUsers[i].Score > scoredUsers[j].Score
	})

	// 4. 选择最优用户
	selectedUsers := make([]int64, 0)
	for i := 0; i < req.Max_users && i < len(scoredUsers); i++ {
		selectedUsers = append(selectedUsers, scoredUsers[i].User.Id)
	}

	// 5. 更新任务负责人
	task.ResponsibleUserIds = selectedUsers
	task.UpdatedAt = time.Now()

	err = l.svcCtx.TaskModel.Update(l.ctx, task)
	if err != nil {
		return nil, err
	}

	// 6. 更新Redis缓存
	taskKey := "task:" + string(rune(task.Id))
	l.svcCtx.Redis.Setex(taskKey, 3600, task)

	// 7. 发布派发事件到MQ
	l.svcCtx.MQ.Publish(l.ctx, "task.dispatch", map[string]interface{}{
		"task_id":            task.Id,
		"company_id":         task.CompanyId,
		"department_id":      task.DepartmentId,
		"responsible_users":  selectedUsers,
		"dispatched_at":      time.Now(),
	})

	// 8. 发送通知邮件给选中的用户
	for _, userId := range selectedUsers {
		user, _ := l.svcCtx.UserAccountModel.FindOne(l.ctx, userId)
		l.svcCtx.Mailer.SendEmail(l.ctx, user.Email, "任务分配通知", 
			fmt.Sprintf("您被分配了新任务：%s", task.Title))
	}

	return &types.ActionResp{Success: true, Message: "任务派发成功"}, nil
}

// 获取候选用户列表
func (l *DispatchLogic) getCandidates(companyId, departmentId int64, requiredTags []string) ([]*core.UserAccount, error) {
	// 这里应该实现根据公司、部门和角色标签查询用户的逻辑
	// 简化实现：查询所有在职用户
	users, err := l.svcCtx.UserAccountModel.FindAll(l.ctx, 1, 1000)
	if err != nil {
		return nil, err
	}

	var candidates []*core.UserAccount
	for _, user := range users {
		if user.Status == 1 && user.CompanyId == companyId {
			candidates = append(candidates, user)
		}
	}

	return candidates, nil
}

// 计算用户得分
func (l *DispatchLogic) calculateScore(user *core.UserAccount, requiredTags []string, companyId, departmentId int64) float64 {
	score := 0.0

	// 技能匹配得分 (40%)
	skillScore := l.calculateSkillMatch(user.RoleTags, requiredTags)
	score += skillScore * 0.4

	// 工作负载得分 (35%)
	workloadScore := l.calculateWorkloadScore(user.Id)
	score += workloadScore * 0.35

	// 部门接近度得分 (25%)
	departmentScore := l.calculateDepartmentScore(user.DepartmentId, departmentId)
	score += departmentScore * 0.25

	return score
}

// 计算技能匹配得分
func (l *DispatchLogic) calculateSkillMatch(userTags []string, requiredTags []string) float64 {
	if len(requiredTags) == 0 {
		return 1.0
	}

	matched := 0
	for _, required := range requiredTags {
		for _, user := range userTags {
			if user == required {
				matched++
				break
			}
		}
	}

	return float64(matched) / float64(len(requiredTags))
}

// 计算工作负载得分（负载越少得分越高）
func (l *DispatchLogic) calculateWorkloadScore(userId int64) float64 {
	// 简化实现：返回随机值
	// 实际应该查询用户当前的任务数量
	return 0.8
}

// 计算部门接近度得分
func (l *DispatchLogic) calculateDepartmentScore(userDeptId, targetDeptId int64) float64 {
	if userDeptId == targetDeptId {
		return 1.0
	}
	// 同公司不同部门
	return 0.5
}

// 得分用户结构
type ScoredUser struct {
	User  *core.UserAccount
	Score float64
}
