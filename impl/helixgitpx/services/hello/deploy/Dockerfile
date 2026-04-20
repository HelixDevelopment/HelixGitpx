# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.work go.work.sum* ./
COPY platform/ ./platform/
COPY gen/ ./gen/
COPY services/hello/ ./services/hello/
RUN cd services/hello && \
    CGO_ENABLED=0 GOWORK=off \
    go build -trimpath -ldflags="-s -w" -o /out/hello ./cmd/hello

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/hello /app/hello
COPY --from=build /src/services/hello/migrations /migrations
USER nonroot
EXPOSE 8001 9001 8081
ENTRYPOINT ["/app/hello"]
LABEL org.opencontainers.image.title="hello"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
LABEL org.opencontainers.image.licenses="Apache-2.0"
