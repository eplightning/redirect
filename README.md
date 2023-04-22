# Redirect

Very simple server serving HTTP redirections.

Useful in scenarios when you're not using any L7 load-balancer and your application doesn't support HTTP->HTTPS redirect.

You can simply add a sidecar to your existing deployment and point plaintext HTTP traffic to the `redirect` container.

## Usage

### Container

Built and ready to use container image is available for `amd64`, `arm64` and `arm` architectures:

```
ghcr.io/eplightning/redirect:v1.0
```

### Configuration

#### Defaults

By default, `redirect` will listen on HTTP connections on port `8080`.

Redirections will use `https` scheme while preserving original host, path and query string. HTTP status code `301 Moved Permanently` is used by default.

#### Environment variables

Behavior can be customized by using environment variables to override default configuration:

```
LISTEN_ADDRESS  - listening address (default :8080)
HOST_OVERRIDE   - host to use for redirects (default use request's Host header)
PATH_OVERRIDE   - path to use for redirects (default use request's path)
QUERY_OVERRIDE  - query to use for redirects (default use request's query string)
SCHEME_OVERRIDE - scheme to use for redirects (default https)
STATUS_CODE     - HTTP status code to use (default 301)
```
