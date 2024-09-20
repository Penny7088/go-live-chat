FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# set the time zone to Shanghai
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY lingua_exchange /app/lingua_exchange
RUN chmod +x /app/lingua_exchange

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./lingua_exchange", "-c", "configs/lingua_exchange.yml"]
# if you use the Configuration Center, lingua_exchange.yml is changed to the Configuration Center configuration.
#CMD ["./lingua_exchange", "-c", "configs/lingua_exchange.yml", "-enable-cc"]
