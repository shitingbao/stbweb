FROM golang

ENV GO111MODULE=on

RUN GOPROXY="https://goproxy.cn" go mod download

COPY . .

WORKDIR /builds/common

RUN go build

EXPOSE 3002

CMD ["common"]