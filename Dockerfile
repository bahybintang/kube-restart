FROM golang:1.18.1 AS builder
WORKDIR /srv/build
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .

FROM alpine:latest
LABEL maintainer="bahybintang@gmail.com"
WORKDIR /root/
COPY --from=builder /srv/build/service .
ENTRYPOINT [ "/root/service" ]
