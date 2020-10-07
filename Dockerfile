FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

COPY . /stbweb

WORKDIR /stbweb/builds/common
RUN go build

EXPOSE 3002

FROM ubuntu

COPY --from=0 /stbweb/builds/common/common .
COPY --from=0 /stbweb/builds/common/config.json .

ENTRYPOINT ["./common"]