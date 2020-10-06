FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io
WORKDIR /stbweb
COPY builds/common /builds/common
COPY core /stbweb
COPY lib /stbweb
COPY loader /stbweb
COPY modules /stbweb

WORKDIR /builds/common
RUN go build

EXPOSE 3002

CMD ["common"]