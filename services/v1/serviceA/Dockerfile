FROM mirror.gcr.io/library/golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/serviceA_v1 ./services/v1/serviceA

FROM mirror.gcr.io/library/alpine:3.19
COPY --from=builder /bin/serviceA_v1 /usr/local/bin/serviceA_v1
EXPOSE 9001
ENTRYPOINT ["/usr/local/bin/serviceA_v1"]
