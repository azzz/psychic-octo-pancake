package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

const contentType = "application/json"

const (
	addItemCommand     = "AddItem"
	removeItemCommand  = "RemoveItem"
	getItemCommand     = "GetItem"
	getAllItemsCommand = "GetAllItems"
)

type Message struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type Client struct {
	queue amqp091.Queue
	conn  *amqp091.Connection
	ch    *amqp091.Channel
}

func (c Client) Close() error {
	var errs []error

	if c.ch != nil {
		if err := c.ch.Close(); err != nil {
			errs = append(errs, fmt.Errorf("cloe channel: %w", err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func New(url, queue string) (*Client, error) {
	var (
		err    error
		client = &Client{}
	)

	client.conn, err = amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	client.ch, err = client.conn.Channel()
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("open channel: %w", err)
	}

	client.queue, err = client.ch.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return client, nil
}

func (c Client) AddItem(ctx context.Context, key, value string) error {
	msg := Message{Command: addItemCommand, Key: key, Value: value}
	return c.send(ctx, msg)
}

func (c Client) RemoveItem(ctx context.Context, key string) error {
	msg := Message{Command: removeItemCommand, Key: key}
	return c.send(ctx, msg)
}

func (c Client) GetItem(ctx context.Context, key string) error {
	msg := Message{Command: getItemCommand, Key: key}
	return c.send(ctx, msg)
}

func (c Client) GetAllItems(ctx context.Context) error {
	msg := Message{Command: getAllItemsCommand}
	return c.send(ctx, msg)
}

func (c Client) send(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	err = c.ch.PublishWithContext(ctx, "", c.queue.Name, false, false, amqp091.Publishing{
		ContentType: contentType,
		Body:        body,
	})

	if err != nil {
		return fmt.Errorf("publish message: %w", err)
	}
	return nil
}
