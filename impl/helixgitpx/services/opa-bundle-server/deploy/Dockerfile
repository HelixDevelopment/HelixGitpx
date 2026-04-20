# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/opa-bundle-server
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/opa-bundle-server ./cmd/opa-bundle-server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/opa-bundle-server /app/opa-bundle-server
USER nonroot
ENTRYPOINT ["/app/opa-bundle-server"]
LABEL org.opencontainers.image.title="opa-bundle-server"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
