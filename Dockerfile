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
		cp ./upsnap ./upsnap-unprivileged &&\
    chmod +x upsnap upsnap-unprivileged &&\
    apk update &&\
    apk add --no-cache libcap &&\
    setcap 'cap_net_raw=+ep' ./upsnap

FROM alpine:3
ARG UPSNAP_HTTP_LISTEN=0.0.0.0:8090
ENV UPSNAP_HTTP_LISTEN=${UPSNAP_HTTP_LISTEN}
RUN apk update &&\
    apk add --no-cache tzdata ca-certificates nmap samba samba-common-tools openssh sshpass curl &&\
    rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=downloader /app/upsnap /app/upsnap-unprivileged ./

RUN printf "%b" '#!'"/bin/sh\n \
pingprivileged=\$(echo \"\$UPSNAP_PING_PRIVILEGED\" | tr '[:upper:]' '[:lower:]') \n \
if [ \"\$pingprivileged\" == \"\" ]; then\n \
  pingprivileged=\"(default) true\"\n \
fi\n \
nonewprivs=\$(cat /proc/self/status | grep \"NoNewPrivs\" | cut -f2)\n \
uid=\$(id -u)\n \
gid=\$(id -g)\n \
netrawcap=\$(cat /proc/self/status | grep \"CapBnd\" | cut -f2)\n \
netrawcap=\${netrawcap:12:1}\n \
ipv4ping_group_range=\$(sysctl net.ipv4.ping_group_range | awk '{print \$3,\$4}')\n \
ipv4ping_group_min=\"\${ipv4ping_group_range% *}\"\n \
ipv4ping_group_max=\"\${ipv4ping_group_range#* }\"\n \
echo \"****************************************************************************\"\n \
echo \"Detected the following:\"\n \
echo \" - NET_RAW capability flag: \$netrawcap\"\n \
echo \" - UID: \$uid, GID: \$gid\"\n \
echo \" - No New Privileges (security-opt: no-new-privileges=true): \$nonewprivs\"\n \
echo \" - UPSNAP_PING_PRIVILEGED: \$pingprivileged\"\n \
echo \" - net.ipv4.ping_group_range: \$ipv4ping_group_range\"\n \
if ! ( [ \"\$pingprivileged\" == \"false\" ] ||\n \
         [ \"\$pingprivileged\" == \"0\" ] ); then\n \
  echo \"Privileged ping configured\"\n \
  if [ \$netrawcap -ne 2 ]; then\n \
    echo \"Unable to start using privileged ping, NET_RAW is not enabled\"\n \
    return 1\n \
  fi\n \
  if [ \$nonewprivs -eq 1 ] && [ \$uid -ne 0 ]; then\n \
    echo \"Unable to start using privileged ping, \\\"no-new-privileges\\\" flag is set and user is not root\"\n \
    return 1\n \
  fi\n \
  echo \"Running Privileged executable\"\n \
  exec /app/upsnap serve --http \${UPSNAP_HTTP_LISTEN}\n \
else\n \
  if [ \"\$ipv4ping_group_range\" == \"1 0\" ] ||\n \
       [ \"\$ipv4ping_group_range\" == \"0 0\" ] ||\n \
       ! ( [ \$gid -ge \$ipv4ping_group_min ] && [ \$gid -le \$ipv4ping_group_max ] ); then\n \
		echo \"Unable to start using unprivileged ping\"\n \
    echo \"Unprivileged ping requires setting \\\"net.ipv4.ping_group_range\\\" correctly for\"\n \
    echo \"Group ID \$gid access. A value of \\\"1 0\\\" means it is disabled.\"\n \
    echo \"Example: sudo sysctl net.ipv4.ping_group_range=\\\"0 2147483647\\\"\"\n \
    return 1\n \
  fi\n \
	echo \"Running Unprivileged executable\"\n \
	exec /app/upsnap-unprivileged serve --http \${UPSNAP_HTTP_LISTEN}\n \
fi\n \
echo \"****************************************************************************\"\n" >/entrypoint.sh && chmod +x /entrypoint.sh


HEALTHCHECK --interval=10s \
    CMD curl -fs "http://${UPSNAP_HTTP_LISTEN}/api/health" || exit 1
ENTRYPOINT ["/entrypoint.sh"]
