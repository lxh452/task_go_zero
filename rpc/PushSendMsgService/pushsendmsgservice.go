package main

import (
	"flag"
	"fmt"

	"task_Project/rpc/PushSendMsgService/internal/config"
	"task_Project/rpc/PushSendMsgService/internal/server"
	"task_Project/rpc/PushSendMsgService/internal/svc"
	"task_Project/rpc/PushSendMsgService/pushSendMsgService"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/pushsendmsgservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pushSendMsgService.RegisterPushSendMsgServiceServer(grpcServer, server.NewPushSendMsgServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
