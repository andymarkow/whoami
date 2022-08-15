# whoami

[![ci](https://github.com/andyglass/whoami/actions/workflows/ci.yml/badge.svg)](https://github.com/andyglass/whoami/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/andyglass/whoami?label=go)
![Docker Pulls](https://img.shields.io/docker/pulls/andyglass/whoami)
![Docker Tag](https://img.shields.io/docker/v/andyglass/whoami?label=docker%20tag)
![Docker Image Size](https://img.shields.io/docker/image-size/andyglass/whoami/latest)


## Routes

<details>
<summary>/ - Get whoami server info</summary>

Request:
```bash
curl -Ss http://localhost | jq
```

Response:
```json
{
  "hostname": "my-laptop",
  "ip": [
    "172.17.0.1",
  ],
  "host": "localhost",
  "url": "/api/v1/status",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/7.59.0"
    ]
  },
  "remote_addr": "172.17.0.1:52550",
  "user_agent": "curl/7.59.0",
  "content_type": "application/json"
}
```
</details>


<details>
<summary>/ping - Check server by ping</summary>

Request:
```bash
curl -Ss http://localhost/ping | jq
```
Response:
```json
{
  "ping": "pong"
}
```
</details>


<details>
<summary>/metrics - Get server metrics in Prometheus format</summary>

Request:
```bash
curl -Ss http://localhost/metrics
```
Response:
```
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
</details>


## Usage

To start docker image run:
```bash
docker run --rm -p 80:80 krezz/whoami
```


## Environment variables

| Environment variable | Default value | Required | Description |
| --- | --- | --- | --- |
| `WEB_SERVER_HOST` | `0.0.0.0` | `false` | Web server listen host |
| `WEB_SERVER_PORT` | `80` | `false` | Web server listen port |