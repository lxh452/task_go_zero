// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handover

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"task_Project/api/task_project/internal/logic/handover"
	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"
)

// 同意交接
func HandoverAcceptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HandoverActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := handover.NewHandoverAcceptLogic(r.Context(), svcCtx)
		resp, err := l.HandoverAccept(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
