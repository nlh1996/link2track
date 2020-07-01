#基础镜像
FROM alpine:latest

WORKDIR /usr/local/src
COPY main /usr/local/src
RUN chmod 777 /usr/local/src/main

ENTRYPOINT ["./main"]