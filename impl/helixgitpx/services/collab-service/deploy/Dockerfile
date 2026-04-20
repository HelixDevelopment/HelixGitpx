# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/collab-service
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/collab-service ./cmd/collab-service

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/collab-service /app/collab-service
USER nonroot
ENTRYPOINT ["/app/collab-service"]
LABEL org.opencontainers.image.title="collab-service"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
