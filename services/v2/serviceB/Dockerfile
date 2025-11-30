FROM mirror.gcr.io/library/golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/serviceB_v2 ./services/v2/serviceB

FROM mirror.gcr.io/library/alpine:3.19
COPY --from=builder /bin/serviceB_v2 /usr/local/bin/serviceB_v2
EXPOSE 9012
ENTRYPOINT ["/usr/local/bin/serviceB_v2"]
