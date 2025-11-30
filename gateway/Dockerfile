FROM mirror.gcr.io/library/golang:1.22-alpine AS builder
WORKDIR /app
# build from repository root to include go.mod
COPY . .
RUN go build -o /bin/lumagate-gateway ./gateway

FROM mirror.gcr.io/library/alpine:3.19
COPY --from=builder /bin/lumagate-gateway /usr/local/bin/lumagate-gateway
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/lumagate-gateway"]
