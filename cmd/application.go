package cmd

import (
	"fmt"
	"github.com/azzz/pillow/internal/client"
	"github.com/rabbitmq/amqp091-go"
)

type Config struct {
	AmqpUrl   string `mapstructure:"amqp_url"`
	AmqpQueue string `mapstructure:"amqp_queue"`
	DataFile  string `mapstructure:"data-file"`
}

func NewClient() (*client.Client, error) {
	conn, err := amqp091.Dial(config.AmqpUrl)
	if err != nil {
		return nil, fmt.Errorf("dial amqp: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open amqp channel: %w", err)
	}

	queue, err := ch.QueueDeclare(config.AmqpQueue, false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("declare amqp queue: %w", err)
	}

	return client.New(queue, ch), nil
}
