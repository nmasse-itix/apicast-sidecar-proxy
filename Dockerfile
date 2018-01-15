# We use the centos base image (and not the "scratch" base image) because we need
# a default trust store (aka SSL/TLS trusted certificates).
# 
# If you want to rebuild this image using the "scratch" base image, make sure to mount a 
# copy of the system truststore at a supported location. 
# 
# See :
#  - https://golang.org/src/crypto/x509/root_unix.go
#  - https://golang.org/src/crypto/x509/root_linux.go
FROM centos:centos7
COPY apicast-sidecar-proxy /apicast-sidecar-proxy
EXPOSE 9090 9091
ENTRYPOINT [ "/apicast-sidecar-proxy" ]
CMD [ ]

