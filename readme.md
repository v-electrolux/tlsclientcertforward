# TLS client certificate forward plugin for traefik

[![Build Status](https://github.com/v-electrolux/tls-client-cert-forward/workflows/Main/badge.svg?branch=main)](https://github.com/v-electrolux/tls-client-cert-forward/actions)

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
--experimental.localplugins.tls-client-cert-forward=true
--experimental.localplugins.tls-client-cert-forward.modulename=github.com/v-electrolux/tls-client-cert-forward
```

- envs as local plugin
```
TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_tls-client-cert-forward=true
TRAEFIK_EXPERIMENTAL_LOCALPLUGINS_tls-client-cert-forward_MODULENAME=github.com/v-electrolux/tls-client-cert-forward
```

- yaml as local plugin
```yaml
experimental:
  localplugins:
    tls-client-cert-forward:
      modulename: github.com/v-electrolux/tls-client-cert-forward
```

- toml as local plugin
```toml
[experimental.localplugins.tls-client-cert-forward]
    modulename = "github.com/v-electrolux/tls-client-cert-forward"
```

### Dynamic config examples

- docker labels
```
traefik.http.middlewares.snForwardMiddleware.plugin.tls-client-cert-forward.snHeaderName=SSL_SN_HEADER
traefik.http.middlewares.snForwardMiddleware.plugin.tls-client-cert-forward.logLevel=warn
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
        tls-client-cert-forward:
          snHeaderName: SSL_SN_HEADER
          logLevel: warn
```
