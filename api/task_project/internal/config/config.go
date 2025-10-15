package config

import "github.com/zeromicro/go-zero/rest"

type RabbitMQConf struct {
	Url       string
	Exchange  string
	Queue     string
	RoutingKey string
}

type SMTPConf struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type MySQLConf struct {
	DataSource string
}

type RedisConf struct {
	Host string
	Type string
	Pass string
	Db   int
}

type Config struct {
	rest.RestConf
    Auth     struct{
        AccessSecret string
        AccessExpire int64
    }
	RabbitMQ RabbitMQConf
	SMTP     SMTPConf
	MySQL    MySQLConf
	Redis    RedisConf
}
