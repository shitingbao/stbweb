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

WORKDIR /stbweb/builds/common

RUN go mod init stbweb
RUN go build

FROM ubuntu
#重新构建，减少体积，这里只需要编译生成的可执行文件，配置文件，前端dist文件即可
COPY --from=stbbuildstage  /stbweb/builds/common/common /opt
COPY --from=stbbuildstage  /stbweb/builds/common/config.json /opt
COPY --from=stbbuildstage  /stbweb/builds/common/dist /opt/dist

EXPOSE 3002

ENTRYPOINT ["./common"]