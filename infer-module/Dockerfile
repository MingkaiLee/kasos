# 第一阶段：编译
# 使用 golang:1.21.6 镜像作为编译环境
FROM golang:1.21.6 AS builder

# 设置工作目录
WORKDIR /app

# 将 Go 项目的文件复制到容器中
COPY . .

# 下载依赖项（如果有）
RUN go mod download

# 编译 Go 程序为静态链接的二进制文件
# -o 指定输出文件名，这里我们命名为 app
# CGO_ENABLED=0 禁用 CGO 以确保生成静态链接的可执行文件
# GOOS=linux 和 GOARCH=amd64 用于跨平台编译（如果需要）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# 第二阶段：创建运行环境
# 自python3.10.13版本开始构建
FROM pyts:v0.1

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件到当前阶段的app目录中
COPY --from=builder /app/app /app

EXPOSE 8080

# 设置容器启动时执行的命令
ENTRYPOINT ["/app/app"]