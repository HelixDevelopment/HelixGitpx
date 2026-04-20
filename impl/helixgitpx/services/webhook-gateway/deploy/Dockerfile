# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/webhook-gateway
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/webhook-gateway ./cmd/webhook-gateway

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/webhook-gateway /app/webhook-gateway
USER nonroot
ENTRYPOINT ["/app/webhook-gateway"]
LABEL org.opencontainers.image.title="webhook-gateway"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
