# services/repo-service/Dockerfile
# Multi-stage, hardened Go service image.
# - Reproducible build (CGO off, stripped, static binary).
# - Distroless runtime; no shell.
# - Non-root user, read-only root filesystem-compatible.
# - Labels for supply-chain (signed + SBOM'd downstream).
# - Version / commit baked in via -ldflags.
#
# Build:
#   docker buildx build --platform=linux/amd64,linux/arm64 \
#     --build-arg VERSION=1.0.0 \
#     --build-arg COMMIT=$(git rev-parse --short HEAD) \
#     -t ghcr.io/vasic-digital/helixgitpx/repo-service:1.0.0 .

# ---------- Build stage --------------------------------------
FROM golang:1.23.2-bookworm AS builder

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOFLAGS="-trimpath -mod=readonly" \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /src

# Pull dep graph first to maximise layer reuse.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

COPY . .

# Build with version metadata; strip debug symbols.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build \
        -ldflags "-s -w \
            -X 'github.com/vasic-digital/helixgitpx/pkg/buildinfo.Version=${VERSION}' \
            -X 'github.com/vasic-digital/helixgitpx/pkg/buildinfo.Commit=${COMMIT}' \
            -X 'github.com/vasic-digital/helixgitpx/pkg/buildinfo.BuildTime=${BUILD_TIME}'" \
        -o /out/repo-service \
        ./services/repo-service/cmd/repo-service

# Verify the binary is static.
RUN file /out/repo-service | grep -q "statically linked"

# ---------- Runtime stage ------------------------------------
# Distroless: no shell, no package manager, no surface area.
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# Standard OCI labels; tooling (syft, trivy) reads these.
LABEL org.opencontainers.image.title="helixgitpx-repo-service"
LABEL org.opencontainers.image.description="HelixGitpx repo-service"
LABEL org.opencontainers.image.vendor="vasic-digital"
LABEL org.opencontainers.image.source="https://github.com/vasic-digital/helixgitpx"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.revision="${COMMIT}"
LABEL org.opencontainers.image.version="${VERSION}"

# Port declaration (informational; actual port binding is at runtime).
EXPOSE 8080 9000 9090 9091

# Run as non-root (uid 65532 from distroless:nonroot).
USER nonroot:nonroot

# Copy the binary.
COPY --from=builder --chown=nonroot:nonroot /out/repo-service /repo-service

# Default health-check (container-level; K8s has its own probes).
HEALTHCHECK --interval=10s --timeout=2s --start-period=5s \
    CMD ["/repo-service", "health", "--address=localhost:9091"]

ENTRYPOINT ["/repo-service"]
CMD ["serve"]
