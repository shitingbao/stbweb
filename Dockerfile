FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /mygo

COPY builds mygo

COPY core mygo

COPY lib mygo

COPY loader mygo

COPY modules mygo

COPY practice mygo

WORKDIR /mygo/builds/common

RUN go build

EXPOSE 3002

CMD ["./common"]