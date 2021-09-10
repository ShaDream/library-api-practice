package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	rootCommand = &cobra.Command{
		Use:   "lib",
		Short: "",
		Long:  "",
		RunE:  nil,
	}

	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file (default is current directory)")

	rootCommand.AddCommand(serveCommand, migrationCommand, feedCommand)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := createConfigFile()
			if err != nil {
				log.Fatalf("Can't create config file: %s", err)
			}
			log.Fatal("Config could not be found. Creating config file in current directory")
		}
		log.Fatalf("Error when read config: %s", err)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

}

func createConfigFile() error {
	viper.Set("db.name", "")
	viper.Set("db.pass", "")
	viper.Set("db.user", "")
	viper.Set("db.host", "")
	viper.Set("db.port", "")
	viper.Set("api.address", "")

	err := viper.SafeWriteConfigAs("./config.toml")
	if err != nil {
		return err
	}
	return nil
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
