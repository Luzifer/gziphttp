ARG SOURCE_DATE_EPOCH=0

FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

ARG SOURCE_DATE_EPOCH
ENV CGO_ENABLED=0

COPY . /src/gziphttp
WORKDIR /src/gziphttp

RUN <<-EOF
  set -ex

  apk add --update git

  install -dm0755 /rootfs/usr/bin

  go build \
    -buildvcs=false \
    -ldflags "-buildid= -s -w -X main.version=$(git describe --tags --always || echo dev)" \
    -mod=readonly \
    -modcacherw \
    -trimpath \
    -o /rootfs/usr/bin/gziphttp

  # Override build / modification time for reproducible build
  find /rootfs -exec touch -d "@${SOURCE_DATE_EPOCH}" {} +
EOF


FROM scratch

LABEL org.opencontainers.image.description="Static file server supporting transparent gzip compression" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/Luzifer/gziphttp" \
      org.opencontainers.image.title="gziphttp" \
      org.opencontainers.image.url="https://github.com/Luzifer/gziphttp" \
      org.opencontainers.image.vendor="Knut Ahlers"

COPY --from=builder /rootfs/ /

EXPOSE 3000
ENTRYPOINT ["/usr/bin/gziphttp"]
