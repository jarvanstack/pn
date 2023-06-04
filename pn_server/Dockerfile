FROM alpine  as runner
WORKDIR /app
# 创建目录
RUN mkdir -p /app
# 解决没有上海时区问题
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai
COPY ./bin/pn_server /app/
COPY ./config-dev.yaml /app/
COPY ./config-prd.yaml /app/
CMD ["./pn_server", "--config", "./config-prd.yaml"]
