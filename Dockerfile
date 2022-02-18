FROM golang:1.13-alpine


ENV GOPROXY="https://goproxy.cn"
ENV CGO_ENABLED=0

RUN apk update && apk add ca-certificates \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf /var/cache/apk/*

WORKDIR $GOPATH/src/jd-sec
ADD . ./
RUN go build -o jdsec -a -installsuffix cgo .
ENTRYPOINT ["./jdsec"]
