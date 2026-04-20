# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/conflict-resolver
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/conflict-resolver ./cmd/conflict-resolver

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/conflict-resolver /app/conflict-resolver
USER nonroot
ENTRYPOINT ["/app/conflict-resolver"]
LABEL org.opencontainers.image.title="conflict-resolver"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
