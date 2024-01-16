FROM --platform=$BUILDPLATFORM golang:1.21-alpine as builder

ARG APP_VERSION
ARG TARGETOS
ARG TARGETARCH

ENV GOOS=$TARGETOS \
    GARCH=$TARGETARCH \
    CGO_ENABLED=0 \
    GO111MODULE=on

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download &&\
    go mod verify

COPY . .
RUN apk update &&\
    apk add --update --no-cache git ca-certificates &&\
    go build -v -a -ldflags "-X main.Version=${APP_VERSION}" -o /build/whoami .


FROM alpine:3.19

COPY --from=builder /build/whoami /usr/local/bin/whoami

CMD ["/usr/local/bin/whoami"]
