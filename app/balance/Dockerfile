FROM alpine:latest
RUN cd /home
WORKDIR /home
COPY balance .
COPY balance.yaml .
EXPOSE 9800
RUN chmod +x balance
ENTRYPOINT ["/home/balance --config=/home/balance.yaml"]