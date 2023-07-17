
FROM golang:alpine

# 开启交叉编译
ENV CGO_ENABLED=1

# 官方并没有提供预编译的包，需要自己编译，但是直接编译的话会提示报错，需要在先安装一下g++
RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++

# 设置工作目录
WORKDIR /app

# 将项目文件复制到容器的工作目录
COPY . .

# 编译 Go 程序
RUN go build -o imaotai

# 设置容器启动时执行的命令
ENTRYPOINT ["./imaotai"]