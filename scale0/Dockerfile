FROM golang:alpine as builder

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/srv ./cmd/srv/...

FROM scratch

WORKDIR /app

COPY --from=builder /app/srv /usr/bin/srv

ENTRYPOINT [ "/usr/bin/srv" ]