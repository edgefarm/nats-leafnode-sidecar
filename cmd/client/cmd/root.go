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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/edgefarm/nats-leafnode-sidecar/pkg/client"
	"github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	natsURI   string
	creds     string
	component string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "Client subscribes and unsubscribes to the corresponding registry.",
	Long: `Client subscribes and unsubscribes to the corresponding registry.
This is supposed to run as a sidecar container and is responsible for
telling the registry about new edgefarm.network credentials for the
component attached to.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("client called")
		registryClient, err := client.NewClient(creds, natsURI, component)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		err = registryClient.Start()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		exit := make(chan bool, 1)
		// Signal handling.
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

			for sig := range c {
				log.Printf("Trapped \"%v\" signal\n", sig)
				switch sig {
				case syscall.SIGINT:
					fmt.Println("1")
					registryClient.Shutdown()
					fmt.Println("2")
					exit <- true
					return
				case syscall.SIGTERM:
					registryClient.Shutdown()
					exit <- true
					return
				}
			}
		}()

		<-exit
		fmt.Println("Goodbye")
		os.Exit(0)

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client.yaml)")
	rootCmd.PersistentFlags().StringVar(&natsURI, "natsuri", "nats://leaf-nats.nats:4222", "Nats URI to connect to, e.g. https://nats.example.com:4222")
	rootCmd.PersistentFlags().StringVar(&creds, "creds", "/nats-credentials", "path where to find the nats credentials")
	rootCmd.PersistentFlags().StringVar(&component, "component", "", "name of the component this is the sidecar for")
	rootCmd.PersistentFlags().StringVar(&common.Remote, "remote", "", "remote Nats URI to connect to, e.g. nats://nats.example.com:4222")
	if common.Remote == "" {
		log.Println("argument 'remote' is required")
		os.Exit(1)
	}
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

		// Search config in home directory with name ".client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".client")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
