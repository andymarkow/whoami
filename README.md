# Whoami Go Web Server

[![ci](https://github.com/andymarkow/whoami/actions/workflows/ci.yml/badge.svg)](https://github.com/andymarkow/whoami/actions/workflows/ci.yml)
[![Go](https://img.shields.io/static/v1?label=go&message=v1.21%2b&color=blue&logo=go)](#)
![License](https://img.shields.io/github/license/andymarkow/whoami)
![Docker Pulls](https://img.shields.io/docker/pulls/andymarkow/whoami)
![Docker Tag](https://img.shields.io/docker/v/andymarkow/whoami?label=docker%20tag)
![Docker Image Size](https://img.shields.io/docker/image-size/andymarkow/whoami/latest)
<!-- ![Release](https://img.shields.io/github/v/release/andymarkow/whoami?display_name=release&include_prereleases&sort=date) -->

Simple Go web server based on net/http library which returns information about web server and HTTP context.

Can be used for development, testing and debugging purposes.


## Install

### Docker

To start web server in docker container run:
```bash
docker run --rm --name=whoami -p 80:8080 andymarkow/whoami
```

Make a HTTP request:
```bash
$ curl -Ssi http://localhost/

HTTP/1.1 200 OK
Content-Length: 334
Content-Type: text/plain; charset=utf-8

RequestID: 5e6309a2-c9e9-41b9-bcbe-ad90afddf475
Hostname: 3405570c9e33
IP: [172.17.0.2]
Host: localhost
URL: /
Method: GET
Proto: HTTP/1.1
Params: map[]
Headers: map[Accept:[*/*] User-Agent:[curl/8.4.0] X-Request-Id:[5e6309a2-c9e9-41b9-bcbe-ad90afddf475]]
UserAgent: curl/8.4.0
RemoteAddr: 10.10.65.1:61337
Environment: map[]
```


### Helm

Under development.


### Configuration

> NOTE: Flags take precedence over environment variables.

| Flag | Environment variable | Default value | Description |
| --- | --- | --- | --- |
| `host` | `WHOAMI_HOST` | `0.0.0.0` | Web server listen address |
| `port` | `WHOAMI_PORT` | `8080` | Web server listen port |
| `log-formatter` | `WHOAMI_LOG_FORMATTER` | `json` | Output log formatter: `fmt` or `json` |
| `log-level` | `WHOAMI_LOG_LEVEL` | `info` | Output log level: `debug`, `info`, `warn`, `error` |
| `access-log` | `WHOAMI_ACCESS_LOG` | `false` | Enable web server access log |
| `access-log-skip-paths` | `WHOAMI_ACCESS_LOG_SKIP_PATHS` | `""` | Comma-separated list of url paths to exclude from access log |
| `read-timeout` | `WHOAMI_READ_TIMEOUT` | `"0s"` | Web server read timeout |
| `read-header-timeout` | `WHOAMI_READ_HEADER_TIMEOUT` | `"0s"` | Web server read header timeout |
| `write-timeout` | `WHOAMI_WRITE_TIMEOUT` | `"0s"` | Web server write timeout |
| `tls-cert` | `WHOAMI_TLS_CERT_FILE` | `""` | TLS certificate file |
| `tls-key` | `WHOAMI_TLS_KEY_FILE` | `""` | TLS private key file |
| `tls-ca` | `WHOAMI_TLS_CA_FILE` | `""` | TLS CA certificate file for mTLS authentication |


## Usage

### HTTP Routes

- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `ANY` | `/` | `?[delay=<duration>]` | Returns web server info in plain text format |

  Parameters:
  - `delay` (Optional): Request delay duration in Go-duration format (ex. 5s, 1m, etc).

  Request:
  ```bash
  curl -Ss http://localhost/
  ```

	Response:
	```bash
	RequestID: 29a11732-d622-4c02-8366-4d834954c5aa
	Hostname: my.local
	IP: 192.168.1.1
	Host: localhost
	URL: /
	Method: GET
	Proto: HTTP/1.1
	UserAgent: PostmanRuntime/7.29.2
	RemoteAddr: [::1]:49995
	```
  ---


- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `ANY` | `/api/*` | `?[delay=<duration>]` | Returns web server info in JSON format |

  Parameters:
  - `delay` (Optional): Request delay duration in Go-duration format (ex. 5s, 1m, etc).

  Request:
  ```bash
  curl -Ss http://localhost/api
  ```

	Response:
	```json
  {
    "request_id": "0d5ab77f-c47c-46f3-88b4-fbd56d854007",
    "hostname": "my.local",
    "host": "localhost",
    "url": "/api",
    "method": "GET",
    "user_agent": "curl/7.79.1",
    "remote_addr": "127.0.0.1:50061",
  }
	```
  ---


- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `ANY` | `/data` | `?[size=<size>]&[unit=<unit>]` | Simulates web server data responce with requested size and unit |

  Parameters:
  - `size` (Optional, default: `1`): Response data size.
  - `unit` (Optional, default: `B`): Response data unit. Possible values: B, KB, MB, GB, TB.

  Request:
  ```bash
  curl -Ss http://localhost/data?size=1&unit=KB
  ```
  ---


- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `POST` | `/upload` | `-` | Simulates web server data upload |

  Request:
  ```bash
  curl -Ss -X POST http://localhost/upload -F file=@file.txt
  ```
  ---


- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `GET` | `/health` | `-` | Returns web server healthcheck status |
  
  Request:
  ```bash
  curl -Ss http://localhost/health
  ```

	Response:
	```json
  {"status": 200}
	```
  ---


- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `POST` | `/health` | `-` | Set web server healthcheck status |
  
  Payload: Valid HTTP status code in range `100..599`.

  Request:
  ```bash
  curl -Ss -X POST -d '418' http://localhost/health
  ```

	Response: Accepted with status `202` on success.
  
  ---


- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `GET` | `/metrics` | `-` | Returns web server metrics in Prometheus format |
  
  Request:
  ```bash
  curl -Ss http://localhost/metrics
  ```

	Response:
	```
  ...
  http_request_duration_seconds_sum{code="200",method="GET",path="/metrics",service=""} 0.018380583
  http_request_duration_seconds_count{code="200",method="GET",path="/metrics",service=""} 2
  whoami_build_info{version="0.0.1"} 1
  whoami_runtime_info{arch="arm64",go_version="go1.21",os="darwin"} 1
  ...
	```
