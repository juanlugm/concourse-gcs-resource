FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gcs ./gcs

FROM alpine:3.18
RUN apk add --no-cache bash jq ca-certificates
COPY --from=builder /out/gcs /opt/resource/gcs
COPY resource /opt/resource/
RUN chmod +x /opt/resource/gcs /opt/resource/*
ENTRYPOINT ["/opt/resource/check"]