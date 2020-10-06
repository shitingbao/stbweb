FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /stbweb
COPY ./builds/common builds/common
COPY ./lib lib
COPY ./loader loader
COPY ./modules modules

WORKDIR /builds/common
RUN go build

EXPOSE 3002

CMD ["./common"]