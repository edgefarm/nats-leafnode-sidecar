/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>
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
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	registry "github.com/edgefarm/nats-leafnode-sidecar/pkg/registry"
)

var (
	cfgFile    string
	natsURI    string
	creds      string
	natsConfig string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "registry",
	Short: "Registry takes care of handling nats leafnode connections",
	Long: `Registry takes care of handling nats leafnode connections. It
	waits for a new client registration which provides all needed information
	to handle this connection. A client can also unregister which results in a
	deletion of this information. Registry modifies the nats configuration file with
	the new information provided to the nats server.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("registry called")
		r, err := registry.NewRegistry(natsConfig, creds, natsURI)
		if err != nil {
			log.Fatal(err)
		}
		err = r.Start()
		if err != nil {
			log.Fatal(err)
		}

		// Signal handling.
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

			for sig := range c {
				switch sig {
				case syscall.SIGINT:
					fallthrough
				case syscall.SIGTERM:
					r.Shutdown()
					return
				}
			}
		}()
		for {
			time.Sleep(time.Second)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.registry.yaml)")
	rootCmd.PersistentFlags().StringVar(&natsURI, "natsuri", "nats://nats.nats:4222", "natsURI to connect to")
	rootCmd.PersistentFlags().StringVar(&natsConfig, "natsconfig", "/config/nats.json", "path to nats config file")
	rootCmd.PersistentFlags().StringVar(&creds, "creds", "/creds", "path where to find the nats credentials")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".registry" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".registry")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
