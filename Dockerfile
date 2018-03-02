FROM golang:1.10 AS builder

WORKDIR /go/src/github.com/sscp/telemetry

COPY . .

RUN make install

FROM debian:stretch

COPY --from=builder /go/bin/telemetry /bin/telemetry
COPY docker_config.yml /docker_config.yml

# .csv and .blog files should be stored in host filesystem, not container
VOLUME /csvs
VOLUME /blogs

# Expose default port to run GRPC service on
EXPOSE 9090

# Expose default port to listen for packets on
EXPOSE 33333/udp

CMD /bin/telemetry -c /docker_config.yml server
