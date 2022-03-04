# ---- Base go + dependencies ----
FROM golang:1.17.2-alpine3.14 AS build

ENV APP_HOME=/tmpdir

WORKDIR $APP_HOME

COPY . $APP_HOME

RUN go build -o bin/frankenstein cmd/frankenstein/main.go

# ---- Final frankenstein image ----
FROM alpine:3.14.2
LABEL Team="Frankenstein" \
    email="xcaballeromartinez@gmail.com"

ENV APP_HOME=/var/www/frankenstein \
    TZ=Europe/Madrid \
    USERNAME=frankenstein \
    UID=1000

RUN adduser -D $USERNAME -u $UID

COPY --from=build --chown=1000  /tmpdir/bin/ $APP_HOME/bin/

USER $USERNAME

WORKDIR $APP_HOME

ENV PATH="${PATH}:${APP_HOME}/bin"
