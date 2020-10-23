FROM golang AS stbbuildstage

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /stbweb

COPY builds/common/main.go builds/common/
COPY builds/common/config.json builds/common/
COPY builds/common/dist builds/common/dist
COPY core core
COPY lib lib
COPY loader loader
COPY modules modules
COPY go.mod go.mod

WORKDIR /stbweb/builds/common

RUN go build

FROM ubuntu
#重新构建，减少体积，这里只需要编译生成的可执行文件，配置文件，前端dist文件即可
WORKDIR /opt
COPY --from=stbbuildstage  /stbweb/builds/common/common .
COPY --from=stbbuildstage  /stbweb/builds/common/config.json .
COPY --from=stbbuildstage  /stbweb/builds/common/dist dist

EXPOSE 3002

ENTRYPOINT ["/opt/common"]