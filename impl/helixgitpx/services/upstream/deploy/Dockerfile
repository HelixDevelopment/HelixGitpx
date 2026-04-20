# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/upstream
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/upstream ./cmd/upstream

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/upstream /app/upstream
USER nonroot
ENTRYPOINT ["/app/upstream"]
LABEL org.opencontainers.image.title="upstream"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
