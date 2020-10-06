FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /mygo

COPY . mygo

WORKDIR builds/common

RUN go build

EXPOSE 3002

CMD ["./common"]