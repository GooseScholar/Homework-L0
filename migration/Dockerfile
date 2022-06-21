FROM postgres:13.7-alpine

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v2.6.0/goose-linux64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /migration
#COPY /migration .

RUN goose postgres "host=${dbHost} user="${dbUsername}" password=${dbPassword} dbname=${dbName} sslmode=${dbsslmode}" up 