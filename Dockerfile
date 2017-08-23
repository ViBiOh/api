FROM scratch

HEALTHCHECK --retries=10 CMD https://localhost:1080/health

COPY script/ca-certificates.crt /etc/ssl/certs/

EXPOSE 1080
ENTRYPOINT [ "/bin/sh" ]

COPY bin/api /bin/sh
