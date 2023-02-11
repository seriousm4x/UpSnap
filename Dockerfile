FROM alpine:3 as downloader
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION
ENV BUILDX_ARCH="${TARGETOS:-linux}_${TARGETARCH:-amd64}${TARGETVARIANT}"
WORKDIR /app
RUN wget https://github.com/seriousm4x/UpSnap/releases/download/${VERSION}/UpSnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    unzip UpSnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    rm -f UpSnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    chmod +x upsnap

FROM alpine:3
RUN apk update &&\
    apk add --no-cache tzdata ca-certificates nmap samba-common-tools openssh sshpass curl &&\
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=downloader /app/upsnap upsnap
HEALTHCHECK --interval=10s \
    CMD curl -fs "http://localhost:8090/api/health" || exit 1
ENTRYPOINT ["./upsnap", "serve", "--http=0.0.0.0:8090"]
