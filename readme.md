# TLS client certificate forward plugin for traefik

[![Build Status](https://github.com/v-electrolux/tlsclientcertforward/workflows/Main/badge.svg?branch=master)](https://github.com/v-electrolux/tlsclientcertforward/actions)

Built-in traefik middleware PassTLSClientCert
let you pass many certificate parameters such as common name or serial number.
But all parameter passed in one header value,
so you can not get just pure certificate serial number in header.
This middleware solves this problem for you. It extracts just pure SN value and put in a header

## Configuration

### Fields meaning
- `snHeaderName`: name of header, in which will be put decimal SN value. 
   Default is Forwarded-Tls-Client-Cert-Dec-Sn
- `logLevel`: `warn`, `info` or `debug`. Default is `info`

### Static config examples

- cli as local plugin
```
--experimental.localplugins.tlsclientcertforward=true
--experimental.localplugins.tlsclientcertforward.modulename=github.com/v-electrolux/tlsclientcertforward
```

- envs as local plugin
```
TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_tlsclientcertforward=true
TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_tlsclientcertforward_MODULENAME=github.com/v-electrolux/tlsclientcertforward
```

- yaml as local plugin
```yaml
experimental:
  localplugins:
     tlsclientcertforward:
      modulename: github.com/v-electrolux/tlsclientcertforward
```

- toml as local plugin
```toml
[experimental.localplugins.tlsclientcertforward]
    modulename = "github.com/v-electrolux/tlsclientcertforward"
```

### Dynamic config examples

- docker labels
```
traefik.http.middlewares.snForwardMiddleware.plugin.tlsclientcertforward.snHeaderName=SSL_SN_HEADER
traefik.http.middlewares.snForwardMiddleware.plugin.tlsclientcertforward.logLevel=warn
traefik.http.routers.snForwardRouter.middlewares=snForwardMiddleware
```

- yaml
```yml
http:

  routers:
    snForwardRouter:
      rule: host(`demo.localhost`)
      service: backend
      entryPoints:
        - web
      middlewares:
        - snForwardMiddleware

  services:
    backend:
      loadBalancer:
        servers:
          - url: 127.0.0.1:5000

  middlewares:
     snForwardMiddleware:
      plugin:
        tlsclientcertforward:
          snHeaderName: SSL_SN_HEADER
          logLevel: warn
```
