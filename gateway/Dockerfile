FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/lumagate-gateway ./

FROM alpine:3.19
COPY --from=builder /bin/lumagate-gateway /usr/local/bin/lumagate-gateway
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/lumagate-gateway"]

