FROM golang:alpine AS builder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/proxy app/proxy/main.go

FROM scratch
ENV TZ Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/proxy /app/proxy
COPY --from=builder /build/app/proxy/proxy.yaml /app/proxy.yaml
CMD ["./proxy"]
