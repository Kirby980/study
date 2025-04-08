# 使用 Ubuntu 22.04（GLIBC 2.35）
FROM ubuntu:22.04

# 或使用 Debian Bookworm（GLIBC 2.36）
FROM debian:bookworm-slim

# 基础配置
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY webook /app/webook
RUN chmod +x /app/webook
ENTRYPOINT ["/app/webook"]