package svc

import (
	"task_Project/api/task_project/internal/config"
    "task_Project/api/task_project/internal/middleware"
    "task_Project/model/core"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
    MQ     MQClient
    Mailer Mailer
    Jwt    *middleware.JwtMiddleware
    DB     sqlx.SqlConn
    Redis  *redis.Redis
    
    // Models
    CompanyModel     core.CompanyModel
    DepartmentModel  core.DepartmentModel
    UserAccountModel core.UserAccountModel
    AuthAccountModel core.AuthAccountModel
    TaskModel        core.TaskModel
    UserTaskLogModel core.UserTaskLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    mq := MustNewMQ(c.RabbitMQ)
    mailer := MustNewMailer(c.SMTP)
    jwt := middleware.NewJwtMiddleware(c.Auth.AccessSecret)
    db := sqlx.NewMysql(c.MySQL.DataSource)
    rds := redis.MustNewRedis(redis.RedisConf{Host: c.Redis.Host, Type: c.Redis.Type, Pass: c.Redis.Pass, Tls: false})
    
    return &ServiceContext{
        Config: c,
        MQ:     mq,
        Mailer: mailer,
        Jwt:    jwt,
        DB:     db,
        Redis:  rds,
        
        // Initialize models
        CompanyModel:     core.NewCompanyModel(db),
        DepartmentModel:  core.NewDepartmentModel(db),
        UserAccountModel: core.NewUserAccountModel(db),
        AuthAccountModel: core.NewAuthAccountModel(db),
        TaskModel:        core.NewTaskModel(db),
        UserTaskLogModel: core.NewUserTaskLogModel(db),
    }
}
