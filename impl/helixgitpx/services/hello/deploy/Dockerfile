# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/hello
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/hello ./cmd/hello

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/hello /app/hello
USER nonroot
ENTRYPOINT ["/app/hello"]
LABEL org.opencontainers.image.title="hello"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
