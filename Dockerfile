# Build
FROM --platform=$BUILDPLATFORM golang:latest AS go_builder
WORKDIR /block

# Cache dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

# Build
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /out/block .

# Copy CA certificates from builder
RUN mkdir -p /out/etc/ssl/certs && \
    cp /etc/ssl/certs/ca-certificates.crt /out/etc/ssl/certs/

# Run
FROM scratch
WORKDIR /app
COPY --from=go_builder /out/block /block
COPY --from=go_builder /out/etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/block" ]