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
	"fmt"
	"net/http"
	"os"
	"strings"
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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the puppet forge webserver",
	Long: `Run this command to start serving your own puppet modules.
You can also enable a fallback proxy to forward the requests to
when you don't have the requested module in your local module
set yet.

You can also enable the caching functionality to speed things up.`,
	Run: func(_ *cobra.Command, _ []string) {
		var err error

		log.Setup(config.Dev)

		config.CacheDir, err = homedir.Expand(config.CacheDir)
		if err != nil {
			log.Log.Fatal(err)
		}
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

		// if set, continuously check modules directory every ModulesScanSec seconds
		// otherwise, check only at startup
		if config.ModulesScanSec > 0 {
			go checkModules(config.ModulesScanSec)
		} else {
			checkModules(config.ModulesScanSec)
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

			if config.FallbackProxyUrl != "" {
				if !config.NoCache {
					if _, err := os.Stat(config.CacheDir); err != nil {
						err = os.MkdirAll(config.CacheDir, os.ModePerm)
						if err != nil {
							log.Log.Fatal(err)
						}
					}
					r.Use(customMiddleware.CacheMiddleware(strings.Split(config.CachePrefixes, ","), config.CacheDir))
				}

				r.Use(customMiddleware.ProxyFallback(config.FallbackProxyUrl, func(status int) bool {
					return status == http.StatusNotFound
				},
					func(r *http.Response, body []byte) {
						if config.ImportProxiedReleases && strings.HasPrefix(r.Request.URL.Path, "/v3/files/") && r.StatusCode == http.StatusOK {
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

			apiRouter := openapi.NewRouter(
				openapi.NewModuleOperationsAPIController(moduleService),
				openapi.NewReleaseOperationsAPIController(releaseService),
				openapi.NewSearchFilterOperationsAPIController(searchFilterService),
				openapi.NewUserOperationsAPIController(userService),
			)

			r.Mount("/", apiRouter)

			log.Log.Infof("Listen on %s:%d", config.Bind, config.Port)
			if config.TlsKeyPath != "" && config.TlsCertPath != "" {
				log.Log.Panic(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Bind, config.Port), config.TlsCertPath, config.TlsKeyPath, r))
			} else {
				log.Log.Panic(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Bind, config.Port), r))
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
	serveCmd.Flags().StringVar(&config.FallbackProxyUrl, "fallback-proxy", "", "optional fallback upstream proxy url")
	serveCmd.Flags().BoolVar(&config.Dev, "dev", false, "enables dev mode")
	serveCmd.Flags().StringVar(&config.CacheDir, "cachedir", "~/.gorge/cache", "cache directory")
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
	var err error

	for {
		err = backend.ConfiguredBackend.LoadModules()
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
