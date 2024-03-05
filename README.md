# ⭐ Gorge

<p align="center">
  <img src="logo.png" width="400" />
  <br />
  <a href="https://github.com/dadav/gorge/releases"><img src="https://img.shields.io/github/release/dadav/gorge.svg" alt="Latest Release"></a>
  <a href="https://pkg.go.dev/github.com/dadav/gorge?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="Go Docs"></a>
  <a href="https://github.com/dadav/gorge/actions"><img src="https://img.shields.io/github/actions/workflow/status/dadav/gorge/build.yml" alt="Build Status"></a>
  <img alt="GitHub License" src="https://img.shields.io/github/license/dadav/gorge">
  <br />
  <br />
Gorge is a go implementation for <a href="https://forgeapi.puppet.com">forgeapi.puppet.com</a>
</p>

## 🅰️ Status

This project is still in an very early stage. Contributions are very welcome.

## 🐶 How it works

You put your modules in the directory `~/.gorge/modules/$module/$release.tar.gz` and gorge will
send them to incoming requests from puppet or r10k.
If the module is not found locally it will forward to request (if configured) to an upstream
forge.
The result of this upstream request will be cached for one day (if not disabled with `--no-cache`).
Usually the request end with the module tarball being downloaded. You can set `--import-proxied-releases`
to automatically import them in your `~/.gorge/modules` directory.

## 🌹 Installation

Via `go install`:

```bash
go install github.com/dadav/gorge@latest
```

## 💎 Usage

```bash
Run this command to start serving your own puppet modules.
You can also enable a fallback proxy to forward the requests to
when you don't have the requested module in your local module
set yet.

You can also enable the caching functionality to speed things up.

Usage:
  gorge serve [flags]

Flags:
      --api-version string        the forge api version to use (default "v3")
      --backend string            backend to use (default "filesystem")
      --bind string               host to listen to (default "127.0.0.1")
      --cache-max-age int         max number of seconds responses should be cached (default 86400)
      --cache-prefixes string     url prefixes to cache (default "/v3/files")
      --cachedir string           cache directory (default "/var/cache/gorge")
      --cors string               allowed cors origins separated by comma (default "*")
      --dev                       enables dev mode
      --fallback-proxy string     optional fallback upstream proxy url
  -h, --help                      help for serve
      --import-proxied-releases   add every proxied modules to local store
      --jwt-secret string         jwt secret (default "changeme")
      --jwt-token-path string     jwt token path (default "~/.gorge/token")
      --modulesdir string         directory containing all the modules (default "/opt/gorge/modules")
      --no-cache                  disables the caching functionality
      --port int                  the port to listen to (default 8080)
      --tls-cert string           path to tls cert file
      --tls-key string            path to tls key file

Global Flags:
      --config string   config file (default is $HOME/.gorge.yaml)
```

## 🐂 Examples

```bash
# just start with defaults
gorge serve

# use fallback forge and cache request of modules and files
gorge serve --fallback-proxy https://forge.puppetlabs.com --cache-prefixes /v3/files,/v3/modules
```

## 🍰 Configuration

You can configure gorge in multiple ways.

Via commandline parameters as seen above.

Via file (`$HOME/.config/gorge.yaml` or `./gorge.yaml`):

```yaml
---
# The forge api version to use. Currently only v3 is supported.
api-version: v3
# The backend type to use. Currently only filesystem is supported.
backend: filesystem
# Max seconds to keep the cached responses.
cache-max-age: 86400
# The host to bind the webservice to.
bind: 127.0.0.1
# The prefixes of requests to cache responses from. Multiple entries must be separated by comma.
cache-prefixes: /v3/files
# The directory to write the cached responses to.
cachedir: ~/.gorge/cache
# Value of the `Access-Control-Allow-Origin` header.
cors: "*"
# Enables the dev mode.
dev: false
# Upstream forge to use when local requests return 404
fallback-proxy:
# Import proxied modules into local backend.
import-proxied-releases: false
# Path to local modules.
modulesdir: ~/.gorge/modules
# Disable cache functionality.
no-cache: false
# Port to bind the webservice to.
port: 8080
# The jwt secret used in the protected endpoint validation
jwt-secret: changeme
# The path to write the jwt token to
jwt-token-path: ~/.gorge/token
# Path to tls cert file
tls-cert: ""
# Path to tls key file
tls-key: ""
```

Via environment:

```bash
GORGE_API_VERSION=v3
GORGE_BACKEND=filesystem
GORGE_BIND=127.0.0.1
GORGE_CACHE_MAX_AGE=86400
GORGE_CACHE_PREFIXES=/v3/files
GORGE_CACHEDIR=~/.gorge/cache
GORGE_CORS="*"
GORGE_DEV=false
GORGE_FALLBACK_PROXY=""
GORGE_IMPORT_PROXIED_RELEASES=false
GORGE_MODULESDIR=~/.gorge/modules
GORGE_NO_CACHE=false
GORGE_PORT=8080
GORGE_JWT_SECRET=changeme
GORGE_JWT_TOKEN_PATH=~/.gorge/token
GORGE_TLS_CERT=""
GORGE_TLS_KEY=""
```

Directories are create automatically and the `~` (tilde) in paths are expanded.

## 🐛 Security

Some endpoints are protected and need a valid jwt token. When gorge first starts,
it will create an admin token in the file `~/.gorge/token`. Use this token
in the Authorization header like this:

`Authorization: Bearer <token>`

In dev mode these security checks are disabled.

## 🐝 Development

The code template for `v3` was generated with this command:

```bash
openapi-generator generate -c config.yaml -g go-server -i forge_api_v3.json -o pkg/gen/v3
```

`go-server` does not support the generation of the auth logic, so I had to create a
[middleware for that](./internal/middleware/auth.go) which uses `jwtauth`.

## 🔑 License

[Apache](./LICENSE)
