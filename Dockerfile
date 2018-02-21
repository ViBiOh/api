FROM scratch

HEALTHCHECK --retries=10 CMD [ "/api", "-url", "https://localhost:1080/health" ]

ENTRYPOINT [ "/api" ]
ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

COPY cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY zoneinfo.zip zoneinfo.zip
COPY bin/api /api
