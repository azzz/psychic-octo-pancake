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
	Run: func(cmd *cobra.Command, args []string) {
		logger.Printf("Using config %s", viper.ConfigFileUsed())
		logger.Printf("amqp-url: %s", config.AmqpUrl)
		logger.Printf("amqp-queue: %s", config.AmqpQueue)
		logger.Printf("data-file: %s", config.DataFile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	logger = log.Default()
}

func initConfig() {
	viper.SetConfigName("pillow")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// write default config file
	if err := viper.ReadInConfig(); err != nil {
		viper.Set("amqp_url", "amqp://root:password@localhost:5672/")
		viper.Set("amqp_queue", "command_requests")
		viper.Set("data-file", path.Join(".", "data.log"))

		cobra.CheckErr(viper.SafeWriteConfig())
	}

	cobra.CheckErr(viper.Unmarshal(&config))
}
