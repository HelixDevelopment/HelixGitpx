# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/repo
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/repo ./cmd/repo

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/repo /app/repo
USER nonroot
ENTRYPOINT ["/app/repo"]
LABEL org.opencontainers.image.title="repo"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
