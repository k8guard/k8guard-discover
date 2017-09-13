FROM varikin/golang-glide-alpine AS build
WORKDIR /go/src/github.com/k8guard/k8guard-discover
COPY ./ ./
RUN apk -U add make
RUN make deps build

FROM alpine
RUN apk -U add ca-certificates
COPY --from=build /go/src/github.com/k8guard/k8guard-discover/k8guard-discover /
EXPOSE 3000
ENTRYPOINT ["/k8guard-discover"]
