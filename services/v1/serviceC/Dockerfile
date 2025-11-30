FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/serviceC_v1 ./services/v1/serviceC

FROM alpine:3.19
COPY --from=builder /bin/serviceC_v1 /usr/local/bin/serviceC_v1
EXPOSE 9003
ENTRYPOINT ["/usr/local/bin/serviceC_v1"]
