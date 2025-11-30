FROM golang:1.22-alpine AS builder
WORKDIR /app
# build from repository root to include go.mod
COPY . .
RUN go build -o /bin/lumagate-gateway ./gateway

FROM alpine:3.19
COPY --from=builder /bin/lumagate-gateway /usr/local/bin/lumagate-gateway
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/lumagate-gateway"]
