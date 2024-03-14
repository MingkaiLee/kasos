# 第一阶段：编译
# 使用 golang:1.21.6 镜像作为编译环境
FROM golang:1.21.6 AS builder

# 设置工作目录
WORKDIR /app

# 将 Go 项目的文件复制到容器中
COPY . .

# 下载依赖项（如果有）
# 如果你使用 Go Modules 管理依赖，请确保在项目根目录下有 go.mod 和 go.sum 文件
RUN go mod download

# 编译 Go 程序为静态链接的二进制文件
# -o 指定输出文件名，这里我们命名为 app
# CGO_ENABLED=0 禁用 CGO 以确保生成静态链接的可执行文件
# GOOS=linux 和 GOARCH=amd64 用于跨平台编译（如果需要）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# 第二阶段：创建运行环境
# 使用 scratch 作为最终镜像的基础镜像
FROM scratch

# 从 builder 阶段复制编译好的二进制文件到当前阶段的根目录
COPY --from=builder /app/app /

EXPOSE 8080

# 设置容器启动时执行的命令
ENTRYPOINT ["/app"]