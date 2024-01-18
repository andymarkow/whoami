variable "DOCKER_IMAGE_REPOSITORY" {
  default = "ghcr.io/andymarkow/whoami"
}

variable "DOCKER_IMAGE_TAG" {
  default = "dev"
}

variable "IS_RELEASE" {
  default = ""
}

group "default" {
  targets = ["whoami"]
}

target "whoami" {
  context    = "."
  dockerfile = "Dockerfile"
  platforms  = [
    "linux/amd64",
    "linux/arm64",
  ]
  args = {
    APP_VERSION = "${DOCKER_IMAGE_TAG}"
  }
  tags = [
    "${DOCKER_IMAGE_REPOSITORY}:${DOCKER_IMAGE_TAG}",
    notequal("", IS_RELEASE) ? "${DOCKER_IMAGE_REPOSITORY}:latest" : "",
  ]
}
