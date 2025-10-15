package logic

import (
	"context"

	"task_Project/rpc/PushSendMsgService/internal/svc"
	"task_Project/rpc/PushSendMsgService/pushSendMsgService"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *pushSendMsgService.Request) (*pushSendMsgService.Response, error) {
	// todo: add your logic here and delete this line

	return &pushSendMsgService.Response{
		Pong: "刘兴洪是帅哥",
	}, nil
}
