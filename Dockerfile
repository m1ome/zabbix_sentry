FROM golang:1.14-alpine as builder

WORKDIR /app
COPY . .

RUN apk add git

# Building main application
RUN set -x \
    && export VERSION=$(git rev-parse --verify HEAD --short) \
    && export LDFLAGS="-w -s -X main.Version=${VERSION}" \
    && export CGO_ENABLED=0 \
    && go build -v -ldflags "${LDFLAGS}" -o /zabbix_sentry .

FROM alpine:3.10

RUN apk add --no-cache ca-certificates

WORKDIR /
COPY --from=builder /zabbix_sentry /zabbix_sentry

CMD ["/zabbix_sentry --help"]