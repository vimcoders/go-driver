FROM alpine:latest
RUN cd /home
WORKDIR /home
COPY echo .
EXPOSE 8080
RUN chmod +x echo
ENTRYPOINT ["/home/echo"]