# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/orgteam
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/orgteam ./cmd/orgteam

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/orgteam /app/orgteam
USER nonroot
ENTRYPOINT ["/app/orgteam"]
LABEL org.opencontainers.image.title="orgteam"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
