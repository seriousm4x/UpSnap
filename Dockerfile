FROM alpine:3 AS downloader
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION
ENV BUILDX_ARCH="${TARGETOS:-linux}_${TARGETARCH:-amd64}${TARGETVARIANT}"
WORKDIR /app
RUN wget https://github.com/seriousm4x/UpSnap/releases/download/${VERSION}/UpSnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    unzip UpSnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    rm -f UpSnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    chmod +x upsnap &&\
    apk update &&\
    apk add --no-cache libcap &&\
    setcap 'cap_net_raw=+ep' ./upsnap

FROM alpine:3
ARG UPSNAP_HTTP_LISTEN=127.0.0.1:8090
ENV UPSNAP_HTTP_LISTEN=${UPSNAP_HTTP_LISTEN}
RUN apk update &&\
    apk add --no-cache tzdata ca-certificates nmap samba samba-common-tools openssh sshpass curl &&\
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=downloader /app/upsnap upsnap
HEALTHCHECK --interval=10s \
    CMD curl -fs "http://${UPSNAP_HTTP_LISTEN}/api/health" || exit 1
ENTRYPOINT ["sh", "-c", "./upsnap serve"]
