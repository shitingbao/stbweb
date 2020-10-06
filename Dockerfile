FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io
WORKDIR /stbweb
COPY . stbweb

WORKDIR /stbweb/builds/common
RUN go build

EXPOSE 3002

CMD ["./common"]