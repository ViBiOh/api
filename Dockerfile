FROM scratch

HEALTHCHECK --retries=10 CMD http://localhost:1080/health

COPY script/ca-certificates.crt /etc/ssl/certs/

ENV ZONEINFO script/zoneinfo.zip
COPY script/zoneinfo.zip script/zoneinfo.zip

EXPOSE 1080
ENTRYPOINT [ "/bin/sh" ]

COPY bin/api /bin/sh
