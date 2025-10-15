// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package notification

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"task_Project/api/task_project/internal/logic/notification"
	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"
)

// 通知已读确认
func NotificationAckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotificationAckReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notification.NewNotificationAckLogic(r.Context(), svcCtx)
		resp, err := l.NotificationAck(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
