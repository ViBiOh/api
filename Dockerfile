FROM scratch

HEALTHCHECK --retries=10 CMD https://localhost:1080/health

ENTRYPOINT [ "/bin/sh" ]
ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY zoneinfo.zip zoneinfo.zip
COPY bin/api /bin/sh
