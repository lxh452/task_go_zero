# 企业任务交接与派发系统数据库设计

## 📋 表结构总览

### 核心业务表
1. **company** - 公司表
2. **department** - 部门表  
3. **position** - 职位表
4. **user** - 用户表（登录认证）
5. **employee** - 员工表（业务信息）
6. **role** - 角色表
7. **employee_role** - 员工角色关联表

### 任务相关表
8. **task** - 任务表
9. **task_node** - 任务节点表
10. **task_log** - 任务日志表
11. **task_handover** - 任务交接表

### 系统功能表
12. **notification** - 通知表
13. **user_permission** - 用户权限表
14. **system_config** - 系统配置表
15. **operation_log** - 操作日志表

## 🔗 表关联关系

### 主要外键关系
```
company (1) ←→ (N) department
company (1) ←→ (N) employee
company (1) ←→ (N) role
company (1) ←→ (N) task

department (1) ←→ (N) position
department (1) ←→ (N) employee
department (1) ←→ (N) task_node
department (1) ←→ (1) department (parent_id)

user (1) ←→ (1) employee
user (1) ←→ (N) user_permission

employee (N) ←→ (N) role (通过 employee_role)
employee (1) ←→ (N) task (creator/assigner)
employee (1) ←→ (N) task_node (executor/leader)
employee (1) ←→ (N) task_log
employee (1) ←→ (N) task_handover
employee (1) ←→ (N) notification

task (1) ←→ (N) task_node
task (1) ←→ (N) task_log
task (1) ←→ (N) task_handover
```

## ⚠️ 重要设计说明

### 1. 用户与员工分离
- **user表**: 纯登录认证信息（用户名、密码、基础个人信息）
- **employee表**: 业务相关信息（工号、部门、职位、技能等）
- 一个用户只能对应一个员工，一个员工只能对应一个用户

### 2. 权限管理
- 基于角色的权限控制（RBAC）
- 支持多角色分配
- 权限可以基于菜单、按钮、接口、数据级别

### 3. 任务管理
- 支持单部门和跨部门任务
- 任务可以分解为多个节点
- 每个节点有独立的执行人和负责人
- 完整的任务交接流程

### 4. 通知系统
- 通知关联到员工而不是用户
- 支持多种通知类型和优先级
- 记录发送者和接收者

### 5. 审计日志
- 完整的操作日志记录
- 支持用户和员工双重关联
- 记录请求参数和响应数据

## 🚀 使用建议

### 创建顺序
```sql
-- 1. 基础表
company → department → position → user → employee

-- 2. 权限表  
role → employee_role

-- 3. 任务表
task → task_node → task_log → task_handover

-- 4. 功能表
notification → user_permission → system_config → operation_log
```

### 索引优化
- 所有外键字段都有索引
- 常用查询字段都有索引
- 软删除字段有索引
- 时间字段有索引

### 数据完整性
- 使用外键约束保证数据一致性
- 软删除设计，保留历史数据
- 级联删除和更新策略合理配置
