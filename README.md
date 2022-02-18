# Packaging and Deployment of Go Microservices to Kubernetes

## How to Run

### Docker Compose

Requirements

- [Docker](https://docs.docker.com/get-docker/)

Deploy

``` bash
docker-compose -f deploy/docker-compose/docker-compose.yaml up
```

Clean up

``` bash
docker-compose -f deploy/docker-compose/docker-compose.yaml down
```

### Kubernetes

Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [K3d](https://k3d.io/v5.2.2/#installation)

Create a cluster

``` bash
grep k3d-registry.localhost /etc/hosts || echo "127.0.0.1 k3d-registry.localhost" | sudo tee -a /etc/hosts
k3d cluster create --config deploy/k8s/k3d.yaml
make -C src/urlshortener/ build
make -C src/urlshortener/ push
```

Deploy using kubectl

``` bash
kubectl apply -f deploy/k8s/specs/ --recursive
```

`or`

Deploy using [Helm](https://helm.sh/docs/intro/install/)

``` bash
helm install urlshortener deploy/k8s/charts/urlshortener
```

Clean up

``` bash
k3d cluster stop dev && k3d cluster delete dev
```

## Test deployment

Wait for the deployment to complete

``` bash
curl localhost:8080/v1/version
curl -X POST -H "Content-Type: application/json" -d '{"url":"https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/"}' localhost:8080/v1/url
```

## Alternative container image building

Requirements

- [ko](https://github.com/google/ko)

Build and push

``` bash
cd src/urlshortener &&
KO_DOCKER_REPO=k3d-registry.localhost:50000 \
PROJECT=github.com/vidbregar/go-microservice \
GIT_TAG=$(git tag --points-at HEAD 2>/dev/null) \
REVISION=$(git rev-parse --short HEAD 2>/dev/null) \
ko publish ./cmd/urlshortener/ &&
cd -
```

Then update deployment's `image:` and
add `command: [ "/ko-app/urlshortener", "-config", "/etc/app/config.yaml", "-secrets", "/etc/app/secrets/" ]` 
