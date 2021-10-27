FROM golang:1.16-alpine as builder
WORKDIR /app
COPY . .

RUN go build && \
    chmod 777 clntgo run.sh

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/clntgo /usr/bin/example-clientgo
COPY --from=builder /app/run.sh /usr/bin/example-run.sh

ENTRYPOINT ["/usr/bin/example-run.sh"]
