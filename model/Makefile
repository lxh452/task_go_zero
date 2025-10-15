
.PHONY: gen-mysql-all-model install-goctl

gen-mysql-all-model:
	goctl model mysql ddl --src "./sql/core.sql" --dir "./core" --style "go-zero"

# 安装 goctl 工具
install-goctl:
	go install github.com/zeromicro/go-zero/tools/goctl@latest


