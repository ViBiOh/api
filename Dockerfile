FROM scratch

HEALTHCHECK --retries=10 CMD http://localhost:1080/health

ENTRYPOINT [ "/bin/sh" ]
ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

COPY zoneinfo.zip zoneinfo.zip

COPY bin/api /bin/sh
