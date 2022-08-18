# whoami

[![ci](https://github.com/andyglass/whoami/actions/workflows/ci.yml/badge.svg)](https://github.com/andyglass/whoami/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/github/go-mod/go-version/andyglass/whoami?label=go)
![Docker Pulls](https://img.shields.io/docker/pulls/andyglass/whoami)
![Docker Tag](https://img.shields.io/docker/v/andyglass/whoami?label=docker%20tag)
![Docker Image Size](https://img.shields.io/docker/image-size/andyglass/whoami/latest)

## Usage

To start docker image run:
```bash
docker run --rm -p 80:80 krezz/whoami
```

## Routes

Follow this docs example:
> Ref: https://www.vaultproject.io/api-docs/secret/kv/kv-v2

## Environment variables

| Environment variable | Default value | Required | Description |
| --- | --- | --- | --- |
| `WEB_SERVER_HOST` | `0.0.0.0` | `false` | Web server listen host |
| `WEB_SERVER_PORT` | `80` | `false` | Web server listen port |
