package cmd

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config Config
	logger *log.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pillow",
	Short: "Pillow is a client-server application where server listens to a queue and client sends commands",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logger = log.Default()
	cobra.OnInitialize(initConfig, printDebug)
}

func printDebug() {
	logger.Printf("DEBUG: Using config %s", viper.ConfigFileUsed())
	logger.Printf("DEBUG: amqp-url: %s", config.AmqpUrl)
	logger.Printf("DEBUG: amqp-queue: %s", config.AmqpQueue)
	logger.Printf("DEBUG: data-file: %s", config.DataFile)
}

func initConfig() {
	viper.SetConfigName("pillow")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// write default config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Printf("INFO: config file not found, create a default one")
		viper.Set("amqp_url", "amqp://root:password@localhost:5672/")
		viper.Set("amqp_queue", "command_requests")
		viper.Set("data-file", path.Join(".", "data.log"))

		cobra.CheckErr(viper.SafeWriteConfig())
	}

	cobra.CheckErr(viper.Unmarshal(&config))
}
