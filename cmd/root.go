package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"yarastore/pkg/config"
)

var conf config.Config

var rootCmd = &cobra.Command{
	Use:   "yarastore",
	Short: "A tool to add and match YARA rules against files",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var configFile string

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file for `yarastore`")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initConfig() {
	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal("Error reading file: ", err)
		}
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("Error Unmarshaling: ", err)
	}
}
