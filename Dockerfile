FROM alpine
ADD k8guard-discover /
EXPOSE 3000
ENTRYPOINT ["/k8guard-discover"]
