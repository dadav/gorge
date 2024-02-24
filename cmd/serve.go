/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

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
	"log"
	"net/http"

	v3 "github.com/dadav/gorge/internal/api/v3"
	openapi "github.com/dadav/gorge/pkg/gen/v3/openapi"
	"github.com/spf13/cobra"
)

var apiVersion string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(_ *cobra.Command, _ []string) {
		if apiVersion == "v3" {
			moduleService := v3.NewModuleOperationsApi()
			releaseService := v3.NewReleaseOperationsApi()
			searchFilterService := v3.NewSearchFilterOperationsApi()
			userService := v3.NewUserOperationsApi()
			handler := openapi.NewRouter(
				openapi.NewModuleOperationsAPIController(moduleService),
				openapi.NewReleaseOperationsAPIController(releaseService),
				openapi.NewSearchFilterOperationsAPIController(searchFilterService),
				openapi.NewUserOperationsAPIController(userService),
			)
			fmt.Println("serve on :8080")
			log.Panic(http.ListenAndServe(":8080", handler))
		} else {
			log.Panicf("%s version not supported", apiVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().StringVar(&apiVersion, "api-version", "v3", "the forge api version to use")
}
