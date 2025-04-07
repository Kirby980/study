# 基础镜像
FROM ubuntu:20.04
# 把编译的打包进来这个镜像,放到工作目录 /app 可以随便更好
COPY webook  /app/webook
WORKDIR /app
#CMD 执行命令
#最佳
ENTRYPOINT  [ "/app/webook" ]