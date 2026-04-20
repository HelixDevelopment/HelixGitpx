# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
WORKDIR /src/services/git-ingress
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -trimpath -ldflags="-s -w" -o /out/git-ingress ./cmd/git-ingress

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/git-ingress /app/git-ingress
USER nonroot
ENTRYPOINT ["/app/git-ingress"]
LABEL org.opencontainers.image.title="git-ingress"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
