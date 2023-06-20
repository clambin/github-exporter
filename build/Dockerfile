FROM alpine

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY bld/github-exporter-$TARGETOS-$TARGETARCH /app/github-exporter

EXPOSE 9090
ENTRYPOINT ["/app/github-exporter"]
CMD []
