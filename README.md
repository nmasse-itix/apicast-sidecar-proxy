# A sidecar for Apicast, enabling proxy support

## Why this project ?

[Apicast](https://github.com/3scale/apicast) is an API gateway that fetches his
configuration and asks for authorization on the 3scale services (hosted in their
cloud). The gateway, however, is deployed on-premise in the customer network.

More and more customers enforce the use of an HTTP proxy for outgoing connections.
Sometimes, an authentication is required to connect to the proxy.

Currently, [apicast does not fully support HTTP proxies](https://issues.jboss.org/browse/THREESCALE-221).
Hence this project provides a sidecar container for Apicast that enables HTTP proxies support.

## Deployment

### Classic install 

First, you will have to build the software. You can do it locally on your Mac using:
```sh
brew install golang
git clone https://github.com/nmasse-itix/apicast-sidecar-proxy.git
GOOS=linux GOARCH=amd64 go build -o apicast-sidecar-proxy src/itix.fr/forward/main.go
```

If you run a RHEL or CentOS server, first install go as explained [here](https://golang.org/doc/install)
and build the software:
```sh
git clone https://github.com/nmasse-itix/apicast-sidecar-proxy.git
go build -o apicast-sidecar-proxy src/itix.fr/forward/main.go
```

You can start the proxy by running:
```sh
export https_proxy="http://my.proxy:1234" # Adjust to your environment
export THREESCALE_PORTAL_ENDPOINT="https://<TENANT>-admin.3scale.net"
export BACKEND_ENDPOINT_OVERRIDE="https://su1.3scale.net"
./apicast-sidecar-proxy
```

The proxy start listening by default on port 9090 (admin portal) and 9091 (backend). It will print on the standard output the requests as they go through the proxy.  

If you need to close the current terminal and leave the proxy running, you can start it this way:
```sh
nohup ./apicast-sidecar-proxy &
```

If needed, you can stop the proxy with:
```
pkill apicast-sidecar-proxy
```

### Docker

You can start the proxy using docker: 
```
docker run -d --name=apicast-sidecar-proxy -p 9090:9090 -p 9091:9091 -e https_proxy=http://my.proxy:1234 -e THREESCALE_PORTAL_ENDPOINT=https://<TENANT>-admin.3scale.net -e BACKEND_ENDPOINT_OVERRIDE=https://su1.3scale.net nmasse/apicast-sidecar-proxy
```

Verify that the container is running with:
```
docker ps
```

You can check the proxy logs using: 

```
docker logs apicast-sidecar-proxy
```

You can stop the proxy with:
```
docker stop apicast-sidecar-proxy
docker rm apicast-sidecar-proxy
```

### OpenShift

TODO

## Check that it works

Check that the communication is working with the following commands:
```
$ curl -D - http://localhost:9090/admin/api/accounts.xml 
HTTP/1.1 403 Forbidden
Content-Type: application/xml; charset=utf-8

<?xml version="1.0" encoding="UTF-8"?><error>Access denied</error>
``` 

```
$ curl -D - http://localhost:9091/transactions/authorize.xml
HTTP/1.1 403 Forbidden
Content-Type: application/vnd.3scale-v2.0+xml

<?xml version="1.0" encoding="UTF-8"?><error code="provider_key_or_service_token_required">Provider key or service token are required</error>
```

## Development

### Build
On Mac, you can use:

```sh
brew install golang
GOOS=linux GOARCH=amd64 go build -o apicast-sidecar-proxy src/itix.fr/forward/main.go
```

### Package

```sh
export VERSION=1.0
export DOCKER_USERNAME="nmasse"
git tag -a v$VERSION -m "Version $VERSION"
docker build -t $DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION .
```

### Pushing your image to DockerHub (Optional)

To push the a new version to Dockerhub:
```sh
docker login -u "$DOCKER_USERNAME" https://index.docker.io/v1/
docker images $DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION --format '{{ .ID }}'
docker tag $(docker images $DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION --format '{{ .ID }}') index.docker.io/$DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION
docker push index.docker.io/$DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION
```

And to make it available by default (`latest` tag).
```sh
docker images $DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION --format '{{ .ID }}'
docker tag $(docker images $DOCKER_USERNAME/apicast-sidecar-proxy:$VERSION --format '{{ .ID }}') index.docker.io/$DOCKER_USERNAME/apicast-sidecar-proxy:latest
docker push index.docker.io/$DOCKER_USERNAME/apicast-sidecar-proxy:latest
```

