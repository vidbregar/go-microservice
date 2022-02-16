# syntax=docker/dockerfile:1

############################
# BUILD
############################
FROM golang:1.17.5-alpine AS builder

ARG PROJECT="github.com/vidbregar/go-microservice"
ARG REVISION
ARG GIT_TAG

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build \
  -ldflags "-s -w \
  -X ${PROJECT}/internal/version.GitTag=${GIT_TAG} \
  -X ${PROJECT}/internal/version.Revision=${REVISION}" \
  -mod=vendor \
  -a -o build/urlshortener \
  ./cmd/urlshortener

############################
# FINAL
############################
FROM gcr.io/distroless/static@sha256:8ad6f3ec70dad966479b9fb48da991138c72ba969859098ec689d1450c2e6c97

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/build/urlshortener /bin/urlshortener

CMD [ "/bin/urlshortener", "-config", "/etc/app/config.yaml", "-secrets", "/etc/app/secrets/" ]
