FROM golang:1.16-alpine as builder
WORKDIR /app
COPY . .

RUN go build && \
    chmod 777 ctrlruntime

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ctrlruntime /usr/bin/example-controllerruntimeapp

ENTRYPOINT ["/usr/bin/example-controllerruntimeapp"]
