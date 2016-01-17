// Copyright © 2016 Hylke Visser
// MIT Licensed - See LICENSE file

package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI Variables
var (
	cfgFile string
)

// RootCmd defines the root of the command tree
var RootCmd = &cobra.Command{
	Use:   "ttntool",
	Short: "The Things Network Toolbox",
	Long:  `Tools for interacting The Things Network.`,
}

// Execute runs on start
func Execute() {
	RootCmd.Execute()
}

// init initializes the configuration and command line flags
func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default \"$HOME/.ttntool.yaml\")")

	RootCmd.PersistentFlags().String("nwkSKey", "2B7E151628AED2A6ABF7158809CF4F3C", "Network Session key")
	RootCmd.PersistentFlags().String("appSKey", "2B7E151628AED2A6ABF7158809CF4F3C", "App Session key")

	RootCmd.PersistentFlags().String("broker", "croft.thethings.girovito.nl", "Broker address")

	viper.BindPFlag("nwkSKey", RootCmd.PersistentFlags().Lookup("nwkSKey"))
	viper.BindPFlag("appSKey", RootCmd.PersistentFlags().Lookup("appSKey"))

	viper.BindPFlag("broker", RootCmd.PersistentFlags().Lookup("broker"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ttntool")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Debug("Using config file")
	}
}
