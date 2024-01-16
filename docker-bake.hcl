variable "DOCKER_IMAGE_REPOSITORY" {
  default = "ghcr.io/andyglass/whoami"
}

variable "DOCKER_IMAGE_TAG" {
  default = "local"
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
  tags = [
    "${DOCKER_IMAGE_REPOSITORY}:${DOCKER_IMAGE_TAG}",
    notequal("", IS_RELEASE) ? "${DOCKER_IMAGE_REPOSITORY}:latest" : "",
  ]
  args = {
    APP_VERSION = "${DOCKER_IMAGE_TAG}"
  }
}
