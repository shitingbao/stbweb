FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

COPY /builds/common /builds/common
COPY /core .
COPY /lib .
COPY /loader .
COPY /modules .

WORKDIR /builds/common

RUN go build

EXPOSE 3002



CMD ["common"]