# Go Whoami Web Server

[![ci](https://github.com/andymarkow/whoami/actions/workflows/ci.yml/badge.svg)](https://github.com/andymarkow/whoami/actions/workflows/ci.yml)
[![Go](https://img.shields.io/static/v1?label=go&message=v1.21%2b&color=blue&logo=go)](#)
![Release](https://img.shields.io/github/v/release/andymarkow/whoami?display_name=release&include_prereleases&sort=date)
![License](https://img.shields.io/github/license/andymarkow/whoami)
![Docker Pulls](https://img.shields.io/docker/pulls/andymarkow/whoami)
![Docker Tag](https://img.shields.io/docker/v/andymarkow/whoami?label=docker%20tag)
![Docker Image Size](https://img.shields.io/docker/image-size/andymarkow/whoami/latest)

Simple Go web server which returns information about web server and HTTP context.
Can be used for development, testing and debugging purposes.


## Install

### Docker

To start docker container run:
```bash
docker run --rm --name=whoami -p 80:8080 andyglass/whoami
```

### Helm

Add Helm chart repository:
```bash
helm repo add andyglass https://andyglass.github.com/helm-charts
```

Install Helm chart:
```bash
helm install whoami andyglass/whoami
```

### Environment variables

| Environment variable | Default value | Required | Description |
| --- | --- | --- | --- |
| `WHOAMI_SERVER_ADDR` | `:8080` | `false` | Web server listen address and port |
| `WHOAMI_LOG_LEVEL` | `info` | `false` | Web server log level. Values: [debug,info,warn,error,fatal] |
| `WHOAMI_LOG_URL_EXCLUDES` | `/favicon.ico,/healthz,/metrics` | `false` | Comma-separated list of urls to exclude from access log |


## Usage

### Query parameters

- Optional param: `delay=0ms`

  Description: Sets time to wait for returning response in Go duration format (eg. 10ms, 1s)

- Optional param: `status=200`

  Description: Sets server response HTTP code (eg. 200, 418)

### Routes

- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `ANY` | `/*` | `?[delay=0ms]&[status=200]` | Get server info in plain text format |

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
  ---

- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `ANY` | `/api*` | `?[delay=0ms]&[status=200]&[pretty]` | Get server info in JSON format (use `pretty` to print formatted) |
  
  Request:
  ```bash
  curl -Ss -X GET localhost/api\?pretty
  ```

	Response:
	```json
  {
    "hostname": "my.local",
    "host": "localhost",
    "url": "/api?pretty",
    "method": "GET",
    ...
    "user_agent": "curl/7.79.1",
    "remote_addr": "127.0.0.1:50061",
    "request_id": "0d5ab77f-c47c-46f3-88b4-fbd56d854007"
  }
	```
  ---

- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `GET` | `/healthz` | `?[delay=0ms]` | Get server healthcheck status |
  
  Request:
  ```bash
  curl -Ss -X GET localhost/healthz
  ```

	Response:
	```json
  {"status":"200"}
	```
  ---

- | Method | Path | Description |
  | --- | --- | --- |
  | `POST` | `/healthz` | Set server healthcheck status |
  
  Payload: Valid HTTP status code in range `200..599`

  Request:
  ```bash
  curl -Ss -X POST -d '418' localhost/healthz
  ```

	Response: NoContent with status `204` on success
  
  ---

- | Method | Path | Params | Description |
  | --- | --- | --- | --- |
  | `GET` | `/version` | `?[delay=0ms]` | Get server version info |
  
  Request:
  ```bash
  curl -Ss -X GET localhost/version
  ```

	Response:
	```json
  {"version":"0.0.1"}
	```
  ---

- | Method | Path | Description |
  | --- | --- | --- |
  | `GET` | `/metrics` | Get server metrics in Prometheus format |
  
  Request:
  ```bash
  curl -Ss -X GET localhost/metrics
  ```

	Response:
	```
  ...
  echo_http_requests_total{method="GET",path="/healthz",status="200"} 2
  echo_http_requests_total{method="POST",path="/healthz",status="204"} 2
  whoami_build_info{version="0.0.1"} 1
  whoami_runtime_info{arch="arm64",go_version="go1.19",os="darwin"} 1
  ...
	```
