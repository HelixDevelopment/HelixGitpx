# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/ai-service
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/ai-service ./cmd/ai-service

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/ai-service /app/ai-service
USER nonroot
ENTRYPOINT ["/app/ai-service"]
LABEL org.opencontainers.image.title="ai-service"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
