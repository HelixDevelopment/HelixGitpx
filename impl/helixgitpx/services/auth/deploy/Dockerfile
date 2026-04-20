# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/auth
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/auth ./cmd/auth

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/auth /app/auth
USER nonroot
ENTRYPOINT ["/app/auth"]
LABEL org.opencontainers.image.title="auth"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
