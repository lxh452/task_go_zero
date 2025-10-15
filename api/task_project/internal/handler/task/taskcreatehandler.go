// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"task_Project/api/task_project/internal/logic/task"
	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"
)

// 创建任务
func TaskCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewTaskCreateLogic(r.Context(), svcCtx)
		resp, err := l.TaskCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
