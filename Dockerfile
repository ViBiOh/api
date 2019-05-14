FROM golang:1.12 as builder

WORKDIR /app
COPY . .

RUN make api \
 && curl -s -o /app/cacert.pem https://curl.haxx.se/ca/cacert.pem \
 && curl -s -o /app/zoneinfo.zip https://raw.githubusercontent.com/golang/go/master/lib/time/zoneinfo.zip

ARG CODECOV_TOKEN
RUN curl -s https://codecov.io/bash | bash

FROM scratch

ENV ZONEINFO zoneinfo.zip
EXPOSE 1080

HEALTHCHECK --retries=10 CMD /api -url http://localhost:1080/health
ENTRYPOINT [ "/api" ]

ARG APP_VERSION
ENV VERSION=${APP_VERSION}

COPY doc /doc
COPY --from=builder /app/cacert.pem /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/zoneinfo.zip /app/bin/api /
