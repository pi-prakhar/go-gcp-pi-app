FROM golang:1.23.1 as builder
WORKDIR /app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY config ./config/
COPY internal ./internal
COPY pkg ./pkg
COPY secrets ./secrets/
 
RUN go mod download

ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -o main ./cmd/user-service

FROM gcr.io/distroless/base
# FROM debian:bullseye-slim

WORKDIR /app
USER 1000:1000
COPY --from=builder --chown=1000:1000  /app/main /app/main
COPY --from=builder --chown=1000:1000  /app/config /app/config
COPY --from=builder --chown=1000:1000  /app/secrets /app/secrets
CMD [ "./main" ]

# Keep the container alive indefinitely
# CMD ["/bin/sh"]