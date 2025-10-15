-- 企业
create table if not exists company (
  id bigint not null auto_increment,
  name varchar(200) not null unique,
  domain varchar(200),
  status smallint not null default 1,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp on update current_timestamp,
  primary key (id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- 部门
create table if not exists department (
  id bigint not null auto_increment,
  company_id bigint not null,
  name varchar(200) not null,
  parent_id bigint null,
  status smallint not null default 1,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp on update current_timestamp,
  unique key uk_company_name (company_id, name),
  key idx_department_company (company_id),
  key idx_department_parent (parent_id),
  primary key (id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- 用户（业务档案）
create table if not exists user_account (
  id bigint not null auto_increment,
  company_id bigint not null,
  department_id bigint null,
  account varchar(120) not null,          -- 登录账号名（业务可见）
  name varchar(120) not null,
  email varchar(200),
  role_tags json not null,
  status smallint not null default 1,     -- 1在职 0离职 2请假
  hired_at date,
  left_at date,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp on update current_timestamp,
  unique key uk_company_account (company_id, account),
  key idx_user_company (company_id),
  key idx_user_dept (department_id),
  primary key (id)
  -- FKs omitted for tool compatibility
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- 登录认证（敏感字段隔离）
create table if not exists auth_account (
  user_id bigint not null,
  password_hash varchar(255) not null,
  last_login_at timestamp null,
  login_failed_count int not null default 0,
  locked_until timestamp null,
  primary key (user_id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- 任务主表
create table if not exists task (
  id bigint not null auto_increment,
  company_id bigint not null,
  department_id bigint null,
  title varchar(300) not null,
  description text,
  attachments_url text,                         -- 外部存储URL（可逗号分隔或JSON清单）
  owner_user_ids json not null,
  node_user_ids json not null,
  flow_assignees json not null,                 -- 节点->用户id[] 映射
  priority smallint not null default 3,         -- 1高 2中 3低
  start_date date,
  due_date date,
  schedule_granularity varchar(10) default 'day', -- day|week|month（前端传值控制展示）
  status smallint not null default 1,           -- 1进行中 0暂停 2完成 3取消
  handover_required tinyint(1) not null default 0,
  created_by bigint null,
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp on update current_timestamp,
  key idx_task_company (company_id),
  key idx_task_dept (department_id),
  key idx_task_due (due_date),
  key idx_task_status (status),
  key idx_task_creator (created_by),
  primary key (id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- 用户任务明细（日更/进度/交接记录）
create table if not exists user_task_log (
  id bigint not null auto_increment,
  company_id bigint not null,
  task_id bigint not null,
  node_key varchar(120),                         -- 对应 flow_assignees 的键（流程节点）
  user_id bigint not null,
  collaborator_user_ids json not null,           -- 协同人（json 数组）
  progress_percent int not null default 0,
  report_note text,                              -- 每日汇报/备注
  status smallint not null default 1,            -- 1进行中 2阻塞 3完成
  deadline_at timestamp null,
  handover_id bigint null,                       -- 交接id（与本表自关联或外部 handover 表，先保留）
  is_active tinyint(1) not null default 1,       -- 当前责任是否有效（交接后旧记录置false）
  created_at timestamp not null default current_timestamp,
  updated_at timestamp not null default current_timestamp on update current_timestamp,
  key idx_log_task (task_id),
  key idx_log_company (company_id),
  key idx_log_user (user_id),
  key idx_log_node (node_key),
  key idx_log_active (is_active),
  primary key (id)
  -- FKs omitted for tool compatibility
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;



