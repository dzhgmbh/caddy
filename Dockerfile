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
    cp caddy /

FROM scratch

EXPOSE 80 443 2015

WORKDIR .caddy

COPY --from=build /caddy /bin/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bin/caddy"]
