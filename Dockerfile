FROM scratch
COPY apicast-sidecar-proxy /apicast-sidecar-proxy
EXPOSE 9090 9091
ENTRYPOINT [ "/apicast-sidecar-proxy" ]

