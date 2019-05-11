FROM golang:1.12 as builder

ARG CODECOV_TOKEN
WORKDIR /app
COPY . .

ENV APP_NAME api

RUN make ${APP_NAME} \
 && curl -s https://codecov.io/bash | bash \
 && curl -s -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && curl -s -o /app/zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

FROM scratch

ARG VERSION
ENV APP_NAME api
ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

HEALTHCHECK --retries=10 CMD [ "/api", "-url", "http://localhost:1080/health" ]
ENTRYPOINT [ "/api" ]

COPY doc /doc
COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/zoneinfo.zip /app/bin/${APP_NAME} /
