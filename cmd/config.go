package cmd

type Config struct {
	AmqpUrl   string `mapstructure:"amqp_url"`
	AmqpQueue string `mapstructure:"amqp_queue"`
	DataFile  string `mapstructure:"data-file"`
}
