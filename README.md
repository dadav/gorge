# ⭐ Gorge

Gorge is a go implementation for [forgeapi.puppet.com](https://forgeapi.puppet.com/).

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
      --api-version string      the forge api version to use (default "v3")
      --backend string          backend to use (default "filesystem")
      --bind string             host to listen to
      --cache-prefixes string   url prefixes to cache (default "/v3/files")
      --cachedir string         cache directory (default "/var/cache/gorge")
      --cors string             allowed cors origins separated by comma (default "*")
      --dev                     enables dev mode
      --fallback-proxy string   optional fallback upstream proxy url
  -h, --help                    help for serve
      --modulesdir string       directory containing all the modules (default "/opt/gorge/modules")
      --no-cache                disables the caching functionality
      --port int                the port to listen to (default 8080)

Global Flags:
      --config string   config file (default is $HOME/.gorge.yaml)
```

## 🐂 Examples

```bash
# use the pupeptlabs forge as fallback
gorge serve --fallback-proxy https://forge.puppetlabs.com

# enable cache for every request
gorge serve --fallback-proxy https://forge.puppetlabs.com --cache-prefixes /v3
```

## 🍰 Configuration

Use the `$HOME/.config/gorge.yaml` (or `./gorge.yaml`):

```yaml
---
api-version: v3
backend: filesystem
bind: 127.0.0.1
cache-prefixes: /v3/files
cachedir: /var/cache/gorge
cors: "*"
dev: false
fallback-proxy:
modulesdir: /opt/gorge/modules
no-cache: false
port: 8080
```

Or the environment:

```bash
GORGE_API_VERSION: v3
GORGE_BACKEND: filesystem
GORGE_BIND: 127.0.0.1
GORGE_CACHE_PREFIXES: /v3/files
GORGE_CACHEDIR: /var/cache/gorge
GORGE_CORS: "*"
GORGE_DEV: false
GORGE_FALLBACK_PROXY:
GORGE_MODULESDIR: /opt/gorge/modules
GORGE_NO_CACHE: false
GORGE_PORT: 8080
```

## 🔑 License

[Apache](./LICENSE)
