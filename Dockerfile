FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

COPY . /stbweb

WORKDIR /stbweb/builds/common
RUN go build

EXPOSE 3002

FROM ubuntu
#重新构建，减少体积，这里只需要编译生成的可执行文件，配置文件，前端dist文件即可
COPY --from=0 /stbweb/builds/common/common .
COPY --from=0 /stbweb/builds/common/config.json .
COPY --from=0 /stbweb/builds/common/dist dist

ENTRYPOINT ["./common"]