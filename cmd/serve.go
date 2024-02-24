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

	v3 "github.com/dadav/gorge/internal/api/v3"
	backend "github.com/dadav/gorge/internal/backend"
	config "github.com/dadav/gorge/internal/config"
	log "github.com/dadav/gorge/internal/log"
	customMiddleware "github.com/dadav/gorge/internal/middleware"
	openapi "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
		log.Setup(config.Dev)

		if config.Backend == "filesystem" {
			backend.ConfiguredBackend = backend.NewFilesystemBackend(config.ModulesDir)
		} else {
			log.Log.Fatalf("Invalid backend: %s", config.Backend)
		}

		backend.ConfiguredBackend.LoadModules()

		if config.ApiVersion == "v3" {
			moduleService := v3.NewModuleOperationsApi()
			releaseService := v3.NewReleaseOperationsApi()
			searchFilterService := v3.NewSearchFilterOperationsApi()
			userService := v3.NewUserOperationsApi()

			r := chi.NewRouter()

			r.Use(middleware.Recoverer)
			r.Use(middleware.RealIP)
			r.Use(customMiddleware.RequireUserAgent)
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:   strings.Split(config.CORSOrigins, ","),
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Content-Type"},
				AllowCredentials: false,
				MaxAge:           300,
			}))

			if !config.NoCache {
				if _, err := os.Stat(config.CacheDir); err != nil {
					err = os.MkdirAll(config.CacheDir, os.ModePerm)
					if err != nil {
						log.Log.Fatal(err)
					}
				}
				r.Use(customMiddleware.CacheMiddleware(strings.Split(config.CachePrefixes, ","), config.CacheDir))
			}

			if config.FallbackProxyUrl != "" {
				r.Use(customMiddleware.ProxyFallback(config.FallbackProxyUrl, func(status int) bool {
					return status == http.StatusNotFound
				}))
			}

			apiRouter := openapi.NewRouter(
				openapi.NewModuleOperationsAPIController(moduleService),
				openapi.NewReleaseOperationsAPIController(releaseService),
				openapi.NewSearchFilterOperationsAPIController(searchFilterService),
				openapi.NewUserOperationsAPIController(userService),
			)

			r.Mount("/", apiRouter)

			log.Log.Infof("Listen on %s:%d", config.Bind, config.Port)
			log.Log.Panic(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Bind, config.Port), r))
		} else {
			log.Log.Panicf("%s version not supported", config.ApiVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&config.ApiVersion, "api-version", "v3", "the forge api version to use")
	serveCmd.Flags().IntVar(&config.Port, "port", 8080, "the port to listen to")
	serveCmd.Flags().StringVar(&config.Bind, "bind", "", "host to listen to")
	serveCmd.Flags().StringVar(&config.ModulesDir, "modulesdir", "/opt/gorge/modules", "directory containing all the modules")
	serveCmd.Flags().StringVar(&config.CacheDir, "cachedir", "/var/cache/gorge", "cache directory")
	serveCmd.Flags().StringVar(&config.CachePrefixes, "cache-prefixes", "/v3/files", "url prefixes to cache")
	serveCmd.Flags().StringVar(&config.Backend, "backend", "filesystem", "backend to use")
	serveCmd.Flags().StringVar(&config.CORSOrigins, "cors", "*", "allowed cors origins separated by comma")
	serveCmd.Flags().StringVar(&config.FallbackProxyUrl, "fallback-proxy", "", "optional fallback upstream proxy url")
	serveCmd.Flags().BoolVar(&config.Dev, "dev", false, "enables dev mode")
	serveCmd.Flags().BoolVar(&config.NoCache, "no-cache", false, "disables the caching functionality")
}
