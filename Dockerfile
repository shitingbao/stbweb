FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

COPY ./builds/common /stbweb/builds/common
COPY ./lib /stbweb/lib
COPY ./loader /stbweb/loader
COPY ./modules /stbweb/modules

WORKDIR /stbweb/builds/common
RUN go build

EXPOSE 3002

CMD ["./common"]