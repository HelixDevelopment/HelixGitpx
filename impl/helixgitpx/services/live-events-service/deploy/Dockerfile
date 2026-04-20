# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/live-events-service
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/live-events-service ./cmd/live-events-service

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/live-events-service /app/live-events-service
USER nonroot
ENTRYPOINT ["/app/live-events-service"]
LABEL org.opencontainers.image.title="live-events-service"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
