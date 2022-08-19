# whoami
---
[![ci](https://github.com/andyglass/whoami/actions/workflows/ci.yml/badge.svg)](https://github.com/andyglass/whoami/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/andyglass/whoami?label=go)
![Docker Pulls](https://img.shields.io/docker/pulls/andyglass/whoami)
![Docker Tag](https://img.shields.io/docker/v/andyglass/whoami?label=docker%20tag)
![Docker Image Size](https://img.shields.io/docker/image-size/andyglass/whoami/latest)

## Install

### Docker

To start docker container run:
```bash
docker run --rm --name=whoami -p 80:8080 andyglass/whoami
```
### Environment variables

| Environment variable | Default value | Required | Description |
| --- | --- | --- | --- |
| `WHOAMI_SERVER_ADDR` | `:8080` | `false` | Web server listen address and port |
| `WHOAMI_LOG_LEVEL` | `info` | `false` | Web server log level. Values: [debug,info,warn,error,fatal] |
| `WHOAMI_LOG_URL_EXCLUDES` | `/favicon.ico,/healthz,/metrics` | `false` | Comma-separated list of urls to exclude from access log |

## Routes

- | Method | Path | Description |
	| --- | --- | --- |
	| `GET` | `/*` | Web server listen address and port |

- Path: `GET /*`

  Request:
  ```bash
  curl -Ss localhost:80
  ```

	Response:
	```bash
	Hostname: my.local
	IP: 192.168.1.1
	Host: localhost:80
	URL: /
	Method: GET
	Proto: HTTP/1.1
	UserAgent: PostmanRuntime/7.29.2
	RemoteAddr: [::1]:49995
	RequestID: 29a11732-d622-4c02-8366-4d834954c5aa
	```
