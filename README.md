# whoami

[![ci](https://github.com/andyglass/whoami/actions/workflows/ci.yml/badge.svg)](https://github.com/andyglass/whoami/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/andyglass/whoami?label=go)
![Docker Pulls](https://img.shields.io/docker/pulls/andyglass/whoami)
![Docker Tag](https://img.shields.io/docker/v/andyglass/whoami?label=docker%20tag)
![Docker Image Size](https://img.shields.io/docker/image-size/andyglass/whoami/latest)

Golang Webserver which returns information about webserver and HTTP context

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

- Path: `GET /*`
  
  Description: Get server info in plain text format

  Request:
  ```bash
  curl -Ss -X GET localhost
  ```

	Response:
	```bash
	Hostname: my.local
	IP: 192.168.1.1
	Host: localhost
	URL: /
	Method: GET
	Proto: HTTP/1.1
	UserAgent: PostmanRuntime/7.29.2
	RemoteAddr: [::1]:49995
	RequestID: 29a11732-d622-4c02-8366-4d834954c5aa
	```

- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `GET` | `/api*` | `[?pretty]` | Get server info in JSON format (use `pretty` to print formatted) |
  
  Request:
  ```bash
  curl -Ss -X GET localhost/api?pretty
  ```

	Response:
	```json
  {
    "hostname": "my.local",
    "ip": [
      "192.168.1.1",
    ],
    "host": "localhost",
    "url": "/api?pretty",
    "params": {
      "pretty": [
        ""
      ]
    },
    "method": "GET",
    "proto": "HTTP/1.1",
    "headers": {
      "Accept": [
        "*/*"
      ],
      "User-Agent": [
        "curl/7.79.1"
      ]
    },
    "user_agent": "curl/7.79.1",
    "remote_addr": "127.0.0.1:50061",
    "request_id": "0d5ab77f-c47c-46f3-88b4-fbd56d854007"
  }
	```
