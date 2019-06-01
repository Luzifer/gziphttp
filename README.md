[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/gziphttp)](https://goreportcard.com/report/github.com/Luzifer/gziphttp)
![](https://badges.fyi/github/license/Luzifer/gziphttp)
![](https://badges.fyi/github/downloads/Luzifer/gziphttp)
![](https://badges.fyi/github/latest-release/Luzifer/gziphttp)

# Luzifer / gziphttp

Very simple HTTP server for serving static files with ability for gzip compression if supported by the client.

## Usage

```console
# gziphttp -h
Usage of gziphttp:
      --listen string      Port/IP to listen on (default ":3000")
      --log-level string   Log level (debug, info, warn, error, fatal) (default "info")
  -d, --serve-dir string   Directory to serve files from (default ".")
      --version            Prints current version and exits
```

Here is an example usage inside a Docker container containing (quite large) compiled JavaScript files:

```dockerfile
FROM golang:alpine as go

RUN set -ex \
 && apk add git \
 && go get -v github.com/Luzifer/gziphttp


FROM node:alpine as node

COPY . /src
WORKDIR /src

RUN set -ex \
 && npm ci \
 && npm run build


FROM alpine:latest

COPY --from=go    /go/bin/gziphttp  /usr/local/bin/
COPY --from=node  /src/dist         /usr/local/share/webtotp

EXPOSE 3000/tcp
CMD ["gziphttp", "-d", "/usr/local/share/webtotp"]
```

In this case `gziphttp` serves compressed files to most (all modern) browsers which ensures the download size of the JavaScript files does not hurt as much as it would without gzip compression.
