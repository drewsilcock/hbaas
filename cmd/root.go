package cmd

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	environment string

	rootCmd = &cobra.Command{
		Use:   "hbaas-server",
		Short: "HBaaS server main command",
		Long:  "Backend server for Happy Birthday as a Service (HBaaS) demo RESTful for birthday greetings.",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"config.toml",
		"Configuration file",
	)
	rootCmd.PersistentFlags().StringVarP(
		&environment,
		"env",
		"e",
		"dev",
		"Environment to run server in. Either 'dev' or 'prod'",
	)
	if err := viper.BindPFlag(
		"env",
		rootCmd.PersistentFlags().Lookup("env"),
	); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Unable to load .env file.")
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("hbaas")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
