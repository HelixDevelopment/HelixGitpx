# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/search-service
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/search-service ./cmd/search-service

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/search-service /app/search-service
USER nonroot
ENTRYPOINT ["/app/search-service"]
LABEL org.opencontainers.image.title="search-service"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
