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

Deploy

``` bash
grep k3d-registry.localhost /etc/hosts || echo "127.0.0.1 k3d-registry.localhost" | sudo tee -a /etc/hosts
k3d cluster create --config deploy/k8s/k3d.yaml
make -C src/urlshortener/ build
make -C src/urlshortener/ push
kubectl apply -f deploy/k8s/specs/ --recursive
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