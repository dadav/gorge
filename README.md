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
If the module is not found locally it will forward the request (if configured) to an upstream
forge.
The results will be cached for one day (if not disabled with `--no-cache`).
Usually the request results in a module tarball being downloaded. You can set `--import-proxied-releases`
to automatically import them in your `~/.gorge/modules` directory.

## 🌹 Installation

Via `go install`:

```bash
go install github.com/dadav/gorge@latest
```

Via `tarball`:

```bash
wget https://github.com/dadav/gorge/releases/download/0.4.2-alpha/gorge_0.4.2-alpha_Linux_x86_64.tar.gz
sudo tar xf gorge_0.4.2-alpha_Linux_x86_64.tar.gz -C /usr/local/bin/ gorge
```

Via `container image`:

```bash
podman run --rm -p 8080:8080 ghcr.io/dadav/gorge:latest
```

Via various package types:

```bash
# rpm
wget https://github.com/dadav/gorge/releases/download/0.4.2-alpha/gorge_0.4.2-alpha_linux_amd64.rpm
sudo yum localinstall gorge_0.4.2-alpha_linux_amd64.rpm

# deb
wget https://github.com/dadav/gorge/releases/download/0.4.2-alpha/gorge_0.4.2-alpha_linux_amd64.deb
sudo apt install gorge_0.4.2-alpha_linux_amd64.deb

# apk
wget https://github.com/dadav/gorge/releases/download/0.4.2-alpha/gorge_0.4.2-alpha_linux_amd64.apk
sudo apk add --allow-untrusted gorge_0.4.2-alpha_linux_amd64.apk
```

Via `helm`:

```bash
git clone https://github.com/dadav/gorge.git
cd gorge/helm/gorge
helm install --namespace gorge --create-namespace gorge .
```

## 💎 Usage

```bash
Run this command to start serving your own puppet modules.
You can also enable one or more fallback proxies to forward the requests to
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
      --cache-by-full-request-uri will cache responses by the full request URI (incl. query fragments) instead of only the request path
      --cors string               allowed cors origins separated by comma (default "*")
      --dev                       enables dev mode
      --drop-privileges           drops privileges to the given user/group
      --fallback-proxy string     optional comma separated list of fallback upstream proxy urls
      --group string              give control to this group or gid (requires root)
  -h, --help                      help for serve
      --import-proxied-releases   add every proxied modules to local store
      --jwt-secret string         jwt secret (default "changeme")
      --jwt-token-path string     jwt token path (default "~/.gorge/token")
      --modules-scan-sec int      seconds between scans of directory containing all the modules. (default 0 means only scan at startup)
      --modulesdir string         directory containing all the modules (default "~/.gorge/modules")
      --no-cache                  disables the caching functionality
      --port int                  the port to listen to (default 8080)
      --tls-cert string           path to tls cert file
      --tls-key string            path to tls key file
      --ui                        enables the web ui
      --user string               give control to this user or uid (requires root)

Global Flags:
      --config string   config file (default is $HOME/.gorge.yaml)
```

### ⛳ Autostart

You can use [gorge.service](./gorge.service) to integrate gorge into your systemd autostart.

The required steps are:

```bash
wget https://raw.githubusercontent.com/dadav/gorge/main/gorge.service
wget https://raw.githubusercontent.com/dadav/gorge/main/defaults.yaml
sudo mv gorge.service /etc/systemd/system/gorge.service
sudo mv defaults.yaml /etc/gorge.yaml
sudo systemctl daemon-reload
sudo systemctl enable --now gorge.service
```

If you've installed gorge as a package (rpm, deb, apk), the required files should already be there.

## 🐂 Examples

```bash
# just start with defaults
gorge serve

# use fallback forge and cache request of modules and files
gorge serve --fallback-proxy https://forge.puppetlabs.com --cache-prefixes /v3/files,/v3/modules

# first use the internal forge server, then (if failed) the official forge and cache request of modules and files
gorge serve --fallback-proxy https://internal-forge.example.com,https://forge.puppetlabs.com --cache-prefixes /v3/files,/v3/modules
```

## 🍰 Configuration

You can configure gorge in multiple ways.

Via commandline parameters as seen above.

Via file (`$HOME/.config/gorge.yaml` or `./gorge.yaml`):

```yaml
---
# Enable basic web ui
ui: false
# Set uid of process to this users uid
user: ""
# Set gid of process to this groups gid
group: ""
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
# Value of the `Access-Control-Allow-Origin` header.
cors: "*"
# Enables the dev mode.
dev: false
# Drop privileges if running as root (user & group options must be set)
drop-privileges: false
# List of comma separated upstream forge(s) to use when local requests return 404
fallback-proxy:
# Import proxied modules into local backend.
import-proxied-releases: false
# Path to local modules.
modulesdir: ~/.gorge/modules
# Seconds between scans of directory containing all the modules
modules-scan-sec: 0
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
GORGE_UI=false
GORGE_USER=""
GORGE_GROUP=""
GORGE_API_VERSION=v3
GORGE_BACKEND=filesystem
GORGE_BIND=127.0.0.1
GORGE_CACHE_MAX_AGE=86400
GORGE_CACHE_PREFIXES=/v3/files
GORGE_CACHE_BY_FULL_REQUEST_URI=false
GORGE_CORS="*"
GORGE_DEV=false
GORGE_DROP_PRIVILEGES=false
GORGE_FALLBACK_PROXY=""
GORGE_IMPORT_PROXIED_RELEASES=false
GORGE_MODULESDIR=~/.gorge/modules
GORGE_MODULES_SCAN_SEC=0
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

### 💊 Using privileged ports (<1024)

If you want to use a port smaller 1024, consider using linux capabilities instead
of running gorge as root.

```bash
# add capability
sudo setcap 'cap_net_bind_service=+ep' /usr/bin/gorge

# run gorge
gorge serve --port 80
```

### 💧 Dropping privileges

There is no need to run gorge as root. But if you still want to do it, be sure to
use the `--drop-privileges` option combined with `--user` and `--group`. You could
set these to `www-data`. It will ensure gorge won't keep running as root, after the
required root actions are done.

```bash
sudo gorge serve --drop-privileges --user www-data --group www-data --port 80
```

## 🐝 Development

The code template for `v3` was generated with this command:

```bash
openapi-generator generate -c config.yaml -g go-server -i forge_api_v3.json -o pkg/gen/v3
```

`go-server` does not support the generation of the auth logic, so I had to create a
[middleware for that](./internal/middleware/auth.go) which uses `jwtauth`.

## 🔑 License

[Apache](./LICENSE)
