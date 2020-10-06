FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

COPY . /stbweb

WORKDIR /stbweb/builds/common
RUN go build

EXPOSE 3002

CMD ["/stbweb/builds/common/common"]