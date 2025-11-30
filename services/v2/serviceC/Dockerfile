FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/serviceC_v2 ./services/v2/serviceC

FROM alpine:3.19
COPY --from=builder /bin/serviceC_v2 /usr/local/bin/serviceC_v2
EXPOSE 9013
ENTRYPOINT ["/usr/local/bin/serviceC_v2"]
