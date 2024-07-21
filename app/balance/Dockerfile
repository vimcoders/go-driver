FROM alpine:latest
# 输出当前用户和主目录
RUN echo "Current user is: $(whoami)" && echo "Home directory is: $(eval echo ~$(whoami))"
WORKDIR /
COPY balance .
EXPOSE 8080
RUN chmod +x balance
ENTRYPOINT ["/balance --config=/balance.yaml"]