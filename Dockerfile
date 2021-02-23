FROM golang:1.15.2

#设置go mod 的代理
RUN export GOPROXY=https://mirrors.aliyun.com/goproxy/
