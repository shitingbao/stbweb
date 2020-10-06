FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

COPY builds/common stbweb/builds/common
COPY core stbweb
COPY lib stbweb
COPY loader stbweb
COPY modules stbweb

WORKDIR /stbweb/builds/common
RUN go build

EXPOSE 3002

CMD ["common"]