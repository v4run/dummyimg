FROM golang:1.15-alpine3.12 as builder
WORKDIR /opt
COPY . /opt
RUN CGO_ENABLED=0 GOOS=linux go build -a -o dummyimg

FROM alpine:3.12
COPY --from=builder /opt/dummyimg /var/www/dummyimg
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /var/www
EXPOSE 80 443
ENTRYPOINT ["/var/www/dummyimg"]
