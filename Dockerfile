FROM golang:1.10.3-alpine as build

RUN apk add --no-cache \
    git

RUN go get \
        github.com/caddyserver/builds

WORKDIR $GOPATH/src/github.com/mholt/caddy

COPY . .

RUN go get -d -v ./caddy/caddymain
RUN cd caddy && \
    go run build.go && \
    cp caddy / && \
    mkdir /.caddy

FROM scratch

EXPOSE 8080 8443 2015

COPY --from=build --chown=65534 /.caddy /.caddy
COPY --from=build /caddy /bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 65534

WORKDIR .caddy

ENTRYPOINT ["/bin/caddy"]
