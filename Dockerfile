FROM golang:1.20 as builder
RUN mkdir /build
WORKDIR /build
COPY . /build
ENV CGO_ENABLED=0
RUN go mod vendor
RUN go build -o k8s-read

FROM scratch

COPY --from=builder /build/k8s-read /k8s-read
COPY --from=builder /build/router.json /router.json
COPY --from=builder /build/cert.pem /cert.pem
COPY --from=builder /build/private.key /private.key

EXPOSE 6100
CMD ["/k8s-read"]

