FROM scratch

HEALTHCHECK --retries=10 CMD http://localhost:1080/health

COPY script/ca-certificates.crt /etc/ssl/certs/
COPY $GOROOT/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip

EXPOSE 1080
ENTRYPOINT [ "/bin/sh" ]

COPY bin/api /bin/sh
