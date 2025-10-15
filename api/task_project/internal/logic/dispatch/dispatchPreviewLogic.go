// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dispatch

import (
	"context"

	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchPreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 自动派发候选预览
func NewDispatchPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchPreviewLogic {
	return &DispatchPreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DispatchPreviewLogic) DispatchPreview(req *types.DispatchPreviewReq) (resp *types.DispatchPreviewResp, err error) {
	// 1. 获取候选用户列表
	candidates, err := l.getCandidates(req.Company_id, req.Department_id, req.Required_role_tags)
	if err != nil {
		return nil, err
	}

	// 2. 计算用户得分并排序
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

	// 3. 构建预览结果
	var candidateList []types.CandidateUser
	for i, scoredUser := range scoredUsers {
		if i >= req.Max_users {
			break
		}
		candidateList = append(candidateList, types.CandidateUser{
			User_id:         scoredUser.User.Id,
			Account:         scoredUser.User.Account,
			Name:            scoredUser.User.Name,
			Email:           scoredUser.User.Email,
			Role_tags:       scoredUser.User.RoleTags,
			Department_id:   scoredUser.User.DepartmentId,
			Score:           scoredUser.Score,
			Skill_match:     l.calculateSkillMatch(scoredUser.User.RoleTags, req.Required_role_tags),
			Workload_score:  l.calculateWorkloadScore(scoredUser.User.Id),
			Department_score: l.calculateDepartmentScore(scoredUser.User.DepartmentId, req.Department_id),
		})
	}

	return &types.DispatchPreviewResp{
		Candidates:     candidateList,
		Total_candidates: int64(len(scoredUsers)),
		Preview_reason: "基于技能匹配、工作负载和部门接近度的派发预览",
		Previewed_at:   time.Now().String(),
	}, nil
}

// 获取候选用户列表
func (l *DispatchPreviewLogic) getCandidates(companyId, departmentId int64, requiredTags []string) ([]*core.UserAccount, error) {
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
func (l *DispatchPreviewLogic) calculateScore(user *core.UserAccount, requiredTags []string, companyId, departmentId int64) float64 {
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
func (l *DispatchPreviewLogic) calculateSkillMatch(userTags []string, requiredTags []string) float64 {
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
func (l *DispatchPreviewLogic) calculateWorkloadScore(userId int64) float64 {
	// 简化实现：返回随机值
	// 实际应该查询用户当前的任务数量
	return 0.8
}

// 计算部门接近度得分
func (l *DispatchPreviewLogic) calculateDepartmentScore(userDeptId, targetDeptId int64) float64 {
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
