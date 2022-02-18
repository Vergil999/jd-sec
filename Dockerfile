FROM golang:1.13-alpine


ENV GOPROXY="https://goproxy.cn"
ENV CGO_ENABLED=0

RUN apk update && apk add ca-certificates \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && rm -rf /var/cache/apk/*

ADD . ./
RUN go build -o jdsec .
RUN chmod 777 start.sh
ENTRYPOINT ["sh","./start.sh"]
