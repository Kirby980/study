.PHONY: docker 
docker:
#  把上次编译的结果删除
	@rm webook || true
#  运行一下 go mod tidy 
	@go mod tidy 
#   指定编译成ARM架构的linux操作系统上可执行的文件 名字叫webook
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=k8s -ldflags="-s -w" -o webook .
#  删除上次编译的docker镜像
	@docker rmi -f Kirby980/webook:v0.0.1 || true
#  构建docker镜像
	@docker build --no-cache -t Kirby980/webook:v0.0.1 .