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
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

const envPrefix = "GORGE"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gorge",
	Short: "Gorge runs a puppet forge server",
	Long:  `You can run this tool to provide access to your puppet modules.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initConfig(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gorge.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	v := viper.New()

	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		homeConfig := filepath.Join(home, ".config")

		// Search config in home directory with name ".gorge" (without extension).
		v.AddConfigPath(homeConfig)
		v.AddConfigPath(".")
		v.SetConfigName("gorge")
		v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		v.SetEnvPrefix(envPrefix)
	}

	v.AutomaticEnv()

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return an error if it's not a missing config file
			return fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	return bindFlags(cmd, v)
}

// bindFlags binds cobra flags with viper config
func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var bindingErrors []string

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				bindingErrors = append(bindingErrors, fmt.Sprintf("failed to bind flag %s: %v", f.Name, err))
			}
		}
	})

	if len(bindingErrors) > 0 {
		return fmt.Errorf("flag binding errors: %s", strings.Join(bindingErrors, "; "))
	}
	return nil
}
