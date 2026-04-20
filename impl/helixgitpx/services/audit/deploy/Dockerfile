# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/audit
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/audit ./cmd/audit

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/audit /app/audit
USER nonroot
ENTRYPOINT ["/app/audit"]
LABEL org.opencontainers.image.title="audit"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
