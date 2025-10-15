// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package company

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"task_Project/api/task_project/internal/logic/company"
	"task_Project/api/task_project/internal/svc"
	"task_Project/api/task_project/internal/types"
)

// 修改企业信息/状态
func CompanyUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompanyUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := company.NewCompanyUpdateLogic(r.Context(), svcCtx)
		resp, err := l.CompanyUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
