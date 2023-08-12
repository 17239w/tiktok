FROM golang:1.20.7 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录为 /app
WORKDIR /app 
# 将当前目录下的所有文件复制到容器的 /app 目录中
COPY . .
RUN go mod tidy
# 切换工作目录到 /app/cmd
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build  -ldflags="-w -s" -o ../main
# 切换工作目录回到 /app
WORKDIR /app
# 创建一个名为 publish 的目录，并将编译后的可执行文件 main 和 conf 目录复制到 publish 目录中
RUN mkdir publish  \
    && cp main publish  \
    && cp -r conf publish

# 指定基础镜像为 Busybox 1.28.4
FROM busybox:1.28.4

# 设置工作目录为 /app
WORKDIR /app

# 从builder阶段的镜像中复制 /app/publish目录到当前镜像的 /app目录
COPY --from=builder /app/publish .

# 指定Gin框架的运行模式为发布模式
ENV GIN_MODE=release
# 容器将监听的端口号为 8080
EXPOSE 8080

# 运行编译后的可执行文件 main
ENTRYPOINT ["./main"]