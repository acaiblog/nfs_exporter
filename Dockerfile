# 使用官方 Go 镜像作为构建环境
FROM golang:1.22.10 AS builder

WORKDIR /go/src/app

# 将本地代码复制到容器中
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nfs_exporter .

# 使用轻量级的 alpine 镜像作为最终镜像
FROM alpine:latest

# 安装必要的依赖（如果需要）
# RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件到最终镜像
COPY --from=builder /go/src/app/nfs_exporter .

# 暴露默认端口
EXPOSE 9689

# 启动命令
ENTRYPOINT ["./nfs_exporter"]