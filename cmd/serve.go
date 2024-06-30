/*
Copyright © 2024 dadav

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	config "github.com/dadav/gorge/internal/config"
	log "github.com/dadav/gorge/internal/log"
	customMiddleware "github.com/dadav/gorge/internal/middleware"
	v3 "github.com/dadav/gorge/internal/v3/api"
	backend "github.com/dadav/gorge/internal/v3/backend"
	openapi "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/stampede"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the puppet forge webserver",
	Long: `Run this command to start serving your own puppet modules.
You can also enable fallback proxies to forward the requests to
when you don't have the requested module in your local module
set yet.

You can also enable the caching functionality to speed things up.`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error

		log.Setup(config.Dev)

		config.ModulesDir, err = homedir.Expand(config.ModulesDir)
		if err != nil {
			log.Log.Fatal(err)
		}
		config.TlsCertPath, err = homedir.Expand(config.TlsCertPath)
		if err != nil {
			log.Log.Fatal(err)
		}
		config.TlsKeyPath, err = homedir.Expand(config.TlsKeyPath)
		if err != nil {
			log.Log.Fatal(err)
		}
		config.JwtTokenPath, err = homedir.Expand(config.JwtTokenPath)
		if err != nil {
			log.Log.Fatal(err)
		}

		if config.Backend == "filesystem" {
			backend.ConfiguredBackend = backend.NewFilesystemBackend(config.ModulesDir)
		} else {
			log.Log.Fatalf("Invalid backend: %s", config.Backend)
		}

		if _, err := os.Stat(config.ModulesDir); err != nil {
			err = os.MkdirAll(config.ModulesDir, os.ModePerm)
			if err != nil {
				log.Log.Fatal(err)
			}
		}

		if config.ApiVersion == "v3" {
			moduleService := v3.NewModuleOperationsApi()
			releaseService := v3.NewReleaseOperationsApi()
			searchFilterService := v3.NewSearchFilterOperationsApi()
			userService := v3.NewUserOperationsApi()

			r := chi.NewRouter()

			// Logger should come before any middleware that modifies the response
			// r.Use(middleware.Logger)
			// Recoverer should also be pretty high in the middleware stack
			r.Use(middleware.Recoverer)
			r.Use(middleware.RealIP)
			r.Use(customMiddleware.RequireUserAgent)
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:   strings.Split(config.CORSOrigins, ","),
				AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
				AllowedHeaders:   []string{"Accept", "Content-Type"},
				AllowCredentials: false,
				MaxAge:           300,
			}))

			if !config.Dev {
				tokenAuth := jwtauth.New("HS256", []byte(config.JwtSecret), nil)
				r.Use(customMiddleware.AuthMiddleware(tokenAuth, func(r *http.Request) bool {
					// Everything but GET is protected and requires a jwt token
					return r.Method != "GET"
				}))

				if _, err = os.Stat(config.JwtTokenPath); err != nil {
					_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user": "admin"})
					err = os.WriteFile(config.JwtTokenPath, []byte(tokenString), 0400)
					if err != nil {
						log.Log.Fatal(err)
					}
					log.Log.Infof("JWT token was written to %s\n", config.JwtTokenPath)
				}
			}

			if !config.NoCache {
				customKeyFunc := func(r *http.Request) uint64 {
					token := r.Header.Get("Authorization")
					return stampede.StringToHash(r.Method, strings.ToLower(token))
				}
				cachedMiddleware := stampede.HandlerWithKey(512, time.Duration(config.CacheMaxAge)*time.Second, customKeyFunc, strings.Split(config.CachePrefixes, ",")...)
				r.Use(cachedMiddleware)
			}

			if config.FallbackProxyUrl != "" {

				proxies := strings.Split(config.FallbackProxyUrl, ",")
				slices.Reverse(proxies)

				for _, proxy := range proxies {
					r.Use(customMiddleware.ProxyFallback(proxy, func(status int) bool {
						return status == http.StatusNotFound
					},
						func(r *http.Response) {
							if config.ImportProxiedReleases && strings.HasPrefix(r.Request.URL.Path, "/v3/files/") && r.StatusCode == http.StatusOK {
								body, err := io.ReadAll(r.Body)
								if err != nil {
									log.Log.Error(err)
									return
								}

								// restore the body
								r.Body = io.NopCloser(bytes.NewBuffer(body))

								release, err := backend.ConfiguredBackend.AddRelease(body)
								if err != nil {
									log.Log.Error(err)
									return
								}
								log.Log.Infof("Imported release %s\n", release.Slug)
							}
						},
					))
				}
			}

			apiRouter := openapi.NewRouter(
				openapi.NewModuleOperationsAPIController(moduleService),
				openapi.NewReleaseOperationsAPIController(releaseService),
				openapi.NewSearchFilterOperationsAPIController(searchFilterService),
				openapi.NewUserOperationsAPIController(userService),
			)

			r.Mount("/", apiRouter)

			r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`{"message": "ok"}`))
			})

			r.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`{"message": "ok"}`))
			})

			bindPort := fmt.Sprintf("%s:%d", config.Bind, config.Port)
			log.Log.Infof("Listen on %s", bindPort)

			ctx, restoreDefaultSignalHandling := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer restoreDefaultSignalHandling()
			g, gCtx := errgroup.WithContext(ctx)

			// if set, continuously check modules directory every ModulesScanSec seconds
			// otherwise, check only at startup
			if config.ModulesScanSec > 0 {
				g.Go(func() error {
					for {
						select {
						case <-gCtx.Done():
							log.Log.Debugln("Canceling module scan goroutine")
							return nil
						case <-time.After(time.Duration(config.ModulesScanSec) * time.Second):
							if err := backend.ConfiguredBackend.LoadModules(); err != nil {
								return err
							}
						}
					}
				})
			} else {
				if err := backend.ConfiguredBackend.LoadModules(); err != nil {
					log.Log.Panic(err)
				}
			}

			server := http.Server{Addr: bindPort, Handler: r, BaseContext: func(_ net.Listener) context.Context { return ctx }}

			g.Go(func() error {
				if config.TlsKeyPath != "" && config.TlsCertPath != "" {
					if err := server.ListenAndServeTLS(config.TlsCertPath, config.TlsKeyPath); err != http.ErrServerClosed {
						return err
					}
				} else {
					if err := server.ListenAndServe(); err != http.ErrServerClosed {
						return err
					}
				}
				return nil
			})

			g.Go(func() error {
				<-gCtx.Done()

				log.Log.Debugln("Shutting down server (timeout: 5s)")
				gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancelShutdown()

				return server.Shutdown(gracefullCtx)
			})

			if err := g.Wait(); err != nil {
				log.Log.Panic(err)
			}
		} else {
			log.Log.Panicf("%s version not supported", config.ApiVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&config.ApiVersion, "api-version", "v3", "the forge api version to use")
	serveCmd.Flags().IntVar(&config.Port, "port", 8080, "the port to listen to")
	serveCmd.Flags().StringVar(&config.Bind, "bind", "127.0.0.1", "host to listen to")
	serveCmd.Flags().StringVar(&config.ModulesDir, "modulesdir", "~/.gorge/modules", "directory containing all the modules")
	serveCmd.Flags().IntVar(&config.ModulesScanSec, "modules-scan-sec", 0, "seconds between scans of directory containing all the modules. (default 0 means only scan at startup)")
	serveCmd.Flags().StringVar(&config.Backend, "backend", "filesystem", "backend to use")
	serveCmd.Flags().StringVar(&config.CORSOrigins, "cors", "*", "allowed cors origins separated by comma")
	serveCmd.Flags().StringVar(&config.FallbackProxyUrl, "fallback-proxy", "", "optional comma separated list of fallback upstream proxy urls")
	serveCmd.Flags().BoolVar(&config.Dev, "dev", false, "enables dev mode")
	serveCmd.Flags().StringVar(&config.CachePrefixes, "cache-prefixes", "/v3/files", "url prefixes to cache")
	serveCmd.Flags().StringVar(&config.JwtSecret, "jwt-secret", "changeme", "jwt secret")
	serveCmd.Flags().StringVar(&config.JwtTokenPath, "jwt-token-path", "~/.gorge/token", "jwt token path")
	serveCmd.Flags().StringVar(&config.TlsCertPath, "tls-cert", "", "path to tls cert file")
	serveCmd.Flags().StringVar(&config.TlsKeyPath, "tls-key", "", "path to tls key file")
	serveCmd.Flags().Int64Var(&config.CacheMaxAge, "cache-max-age", 86400, "max number of seconds responses should be cached")
	serveCmd.Flags().BoolVar(&config.NoCache, "no-cache", false, "disables the caching functionality")
	serveCmd.Flags().BoolVar(&config.ImportProxiedReleases, "import-proxied-releases", false, "add every proxied modules to local store")
}

func checkModules(sleepSeconds int) {
	for {
		err := backend.ConfiguredBackend.LoadModules()
		if err != nil {
			log.Log.Fatal(err)
		}
		if sleepSeconds > 0 {
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
		} else {
			break
		}
	}
}
