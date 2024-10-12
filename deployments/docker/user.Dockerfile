FROM golang:1.23.1 as builder
WORKDIR /app

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY config ./config
COPY internal ./internal
COPY pkg ./pkg
COPY secrets ./secrets
 
RUN go mod download

ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -o main ./cmd/user-service

FROM gcr.io/distroless/base
WORKDIR /app
USER 1000:1000
COPY --from=builder --chown=1000:1000  /app/main /app/main
COPY --from=builder --chown=1000:1000  /app/config /app/config
COPY --from=builder --chown=1000:1000  /app/secrets /app/secrets
ENV CONFIG_FILE_PATH=./config/user-service-config.yaml
ENV GOOGLE_APPLICATION_CREDENTIALS="./secrets/poetic-orb-430304-r9-f97ac87421d0.json"
CMD [ "./main" ]
