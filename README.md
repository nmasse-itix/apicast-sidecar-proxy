# A sidecar for Apicast, enabling proxy support

## Why this project ?

[Apicast](https://github.com/3scale/apicast) is an API gateway that fetches his
configuration and asks for authorization on the 3scale services (hosted in their
cloud). The gateway, however, is deployed on-premise in the customer network.

More and more customers enforce the use of an HTTP proxy for outgoing connections.
Sometimes, an authentication is required to connect to the proxy.

Currently, apicast does not fully support HTTP proxies. Hence this project provides
a sidecar container for Apicast that enables HTTP proxies support.

## Deployment

TODO

## Development

### Build
```
GOOS=linux GOARCH=amd64 go build -o apicast-sidecar-proxy src/itix.fr/forward/main.go
```

### Package

```
VERSION=1.0
git tag TODO
docker build -t apicast-sidecar-proxy:$VERSION .
```

### Pushing your image to DockerHub (Optional)

```
docker login https://index.docker.io/v1/
docker images apicast-sidecar-proxy:$VERSION --format '{{ .ID }}'
docker tag $(docker images apicast-sidecar-proxy:$VERSION --format '{{ .ID }}') index.docker.io/<your-username>/apicast-sidecar-proxy:$VERSION
docker push index.docker.io/<your-username>/apicast-sidecar-proxy:$VERSION
```
