# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/adapter-pool
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/adapter-pool ./cmd/adapter-pool

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/adapter-pool /app/adapter-pool
USER nonroot
ENTRYPOINT ["/app/adapter-pool"]
LABEL org.opencontainers.image.title="adapter-pool"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
