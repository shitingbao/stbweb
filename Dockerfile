FROM golang

ENV GO111MODULE=on

ENV GOPROXY=https://goproxy.io

WORKDIR /builds/common

COPY /builds/common .

COPY /core /

COPY /lib /

COPY /loader /

COPY /modules /

COPY /practice /

RUN go build

EXPOSE 3002

CMD ["common"]