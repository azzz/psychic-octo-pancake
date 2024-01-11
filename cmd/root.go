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
	Use:   "psychic-octo-pancake",
	Short: "psychic-octo-pancake is a consumer/produces application with a queue",
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
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// write default config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Printf("INFO: config file not found, create a default one")

		// I suppose it's fine to pre-hardcode default AMQP URL and Queue name for a while
		// instead of running a sort of initializing wizzard.\
		// RabbitMQ URL
		viper.Set("amqp_url", "amqp://guest:guest@localhost:5672/")
		// Queue name to send commands to a server
		viper.Set("amqp_queue", "command_requests")
		// Server data file to write commands output
		viper.Set("data-file", path.Join(".", "data.log"))

		cobra.CheckErr(viper.SafeWriteConfig())
	}

	cobra.CheckErr(viper.Unmarshal(&config))
}
