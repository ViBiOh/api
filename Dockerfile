FROM golang:1.10 as builder

ENV APP_NAME api
ENV WORKDIR ${GOPATH}/src/github.com/ViBiOh/go-api

WORKDIR ${WORKDIR}
COPY ./ ${WORKDIR}/

RUN make ${APP_NAME} \
 && mkdir -p /app \
 && curl -s -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && curl -s -o /app/zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip \
 && cp bin/${APP_NAME} /app/

FROM scratch

ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

HEALTHCHECK --retries=10 CMD [ "/api", "-url", "https://localhost:1080/health" ]
ENTRYPOINT [ "/api" ]

COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/ /
