# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/sync-orchestrator
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/sync-orchestrator ./cmd/sync-orchestrator

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/sync-orchestrator /app/sync-orchestrator
USER nonroot
ENTRYPOINT ["/app/sync-orchestrator"]
LABEL org.opencontainers.image.title="sync-orchestrator"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
