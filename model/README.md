# Model层使用说明

## 📁 文件结构
```
model/
├── sql/                    # SQL文件目录
│   ├── user.sql           # 用户和员工表
│   ├── company.sql        # 公司、部门、职位表
│   ├── role.sql           # 角色、权限、配置、日志表
│   ├── task.sql           # 任务相关表
│   ├── user_auth.sql      # 用户权限和通知表
│   ├── init.sql           # 完整初始化脚本（包含外键）
│   └── README.md          # 数据库设计文档
├── Makefile               # 构建脚本
└── README.md              # 本文件
```

## 🚀 使用方法

### 1. 单独生成某个表
```bash
# 生成用户表
make user

# 生成公司表
make company

# 生成角色表
make role

# 生成权限表
make user_auth

# 生成任务表
make task
```

### 2. 生成所有表
```bash
# 按依赖顺序生成所有表
make all

# 或者
make
```

### 3. 清理和重建
```bash
# 清理所有生成的文件
make clean

# 清理后重新生成所有表
make rebuild
```

## ⚠️ 重要说明

### 生成顺序
由于表之间存在依赖关系，必须按照以下顺序生成：
1. **user** - 用户表（无依赖）
2. **company** - 公司表（依赖user）
3. **role** - 角色表（依赖company和user）
4. **user_auth** - 权限表（依赖user）
5. **task** - 任务表（依赖company和user）

### 外键约束
- 各个独立的SQL文件中**不包含外键约束**
- 外键约束统一在 `init.sql` 文件中添加
- 这样可以避免表之间的强耦合，支持单独生成

### 数据库初始化
如果需要创建完整的数据库（包含外键约束），请使用：
```sql
-- 执行 init.sql 文件
source model/sql/init.sql;
```

## 🔧 故障排除

### 1. 生成失败
如果某个表生成失败，请检查：
- 是否按照正确的依赖顺序生成
- SQL语法是否正确
- 是否有语法错误

### 2. 外键错误
如果遇到外键相关错误：
- 确保先创建被引用的表
- 检查外键名称是否冲突
- 使用 `init.sql` 进行完整初始化

### 3. 清理重建
如果遇到奇怪的问题：
```bash
make clean
make rebuild
```

## 📊 表依赖关系图
```
user (基础表)
├── company (依赖user)
│   ├── department (依赖company)
│   │   └── position (依赖department)
│   └── employee (依赖user, company, department, position)
├── role (依赖company)
│   └── employee_role (依赖employee, role)
├── user_permission (依赖user)
├── notification (依赖employee)
└── operation_log (依赖user, employee)

task (依赖company, employee)
├── task_node (依赖task, department, employee)
├── task_log (依赖task, task_node, employee)
└── task_handover (依赖task, employee)
```
