FROM alpine:3 as downloader
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION
ENV BUILDX_ARCH="${TARGETOS:-linux}_${TARGETARCH:-amd64}${TARGETVARIANT}"
WORKDIR /app
RUN wget https://github.com/seriousm4x/upsnap/releases/download/${VERSION}/upsnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    unzip upsnap_${VERSION}_${BUILDX_ARCH}.zip &&\
    chmod +x /upsnap

FROM alpine:3
RUN apk update &&\
    apk add --no-cache ca-certificates nmap samba-common-tools &&\
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=downloader /app/upsnap upsnap
ENTRYPOINT ["./upsnap", "serve", "--http=0.0.0.0:8090"]
