FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/serviceA_v2 ./services/v2/serviceA

FROM alpine:3.19
COPY --from=builder /bin/serviceA_v2 /usr/local/bin/serviceA_v2
EXPOSE 9011
ENTRYPOINT ["/usr/local/bin/serviceA_v2"]
