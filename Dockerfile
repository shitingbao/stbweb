FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /stbweb

COPY . . 

RUN cd builds/common

RUN go build

EXPOSE 3002

CMD ["./common"]