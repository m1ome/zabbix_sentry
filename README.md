# Zabbix Sentry
> Zabbix + Sentry = ❤️

## Usage
```bash
Usage of zabbix_sentry:
  -host string
        hostname will be sent to zabbix (default "zabbix-sentry")
  -projects string
        projects to be filtered with, comma separated strings
  -sentry-api-key string
        Sentry api key
  -sentry-url string
        Sentry custom entrypoint (default "http://sentry.io/api/0/")
  -verbose
        verbose mode
  -zabbix-host string
        Zabbix host (default "127.0.0.1")
  -zabbix-port int
        Zabbix port (default 10051)
```

## Events
Events will be thrown in following format:
```
sentry.<org>.<project>.events.<status>
```

## Starting with docker
```bash
docker run --rm -ti w1n2k/zabbix-sentry zabbix_sentry --help
```