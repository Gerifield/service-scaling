FROM golang:alpine as builder

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/srv ./cmd/srv/...

FROM gcr.io/distroless/static as final

WORKDIR /app

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/srv /usr/bin/srv

ENTRYPOINT [ "/usr/bin/srv" ]