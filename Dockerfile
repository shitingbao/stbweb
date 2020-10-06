FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /builds/common

COPY . .

RUN go build

EXPOSE 3002

CMD ["./common"]