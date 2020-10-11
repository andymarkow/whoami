FROM golang:1.15-alpine as builder

ENV GOOS=linux \
    GARCH=amd64 \
    CGO_ENABLED=0 \
    GO111MODULE=on

WORKDIR /workspace

COPY . .

RUN apk update && \
    apk add --update --no-cache ca-certificates && \
    go mod download && \
    go mod verify && \
    go build -x -v -a  -o /build/app .


FROM alpine:latest

COPY --from=builder /build/app /usr/local/bin

EXPOSE 80
CMD ["app"]