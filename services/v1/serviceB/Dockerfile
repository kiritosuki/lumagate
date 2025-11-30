FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/serviceB_v1 ./

FROM alpine:3.19
COPY --from=builder /bin/serviceB_v1 /usr/local/bin/serviceB_v1
EXPOSE 9002
ENTRYPOINT ["/usr/local/bin/serviceB_v1"]

