# syntax=docker/dockerfile:1.7
FROM golang:1.23-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.work go.work.sum* ./
COPY platform/ ./platform/
COPY gen/ ./gen/
COPY services/audit/ ./services/audit/
RUN cd services/audit && \
    CGO_ENABLED=0 GOWORK=off \
    go build -trimpath -ldflags="-s -w" -o /out/audit ./cmd/audit && \
    CGO_ENABLED=0 GOWORK=off \
    go build -trimpath -ldflags="-s -w" -o /out/audit-merkle ./cmd/audit-merkle

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/audit /app/audit
COPY --from=build /out/audit-merkle /app/audit-merkle
COPY --from=build /src/services/audit/migrations /migrations
USER nonroot
EXPOSE 9004 8084
ENTRYPOINT ["/app/audit"]
LABEL org.opencontainers.image.title="audit"
LABEL org.opencontainers.image.source="https://github.com/helixgitpx/helixgitpx"
LABEL org.opencontainers.image.licenses="Apache-2.0"
