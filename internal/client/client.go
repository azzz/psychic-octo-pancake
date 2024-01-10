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
	ch    *amqp091.Channel
}

func New(queue amqp091.Queue, ch *amqp091.Channel) *Client {
	return &Client{queue: queue, ch: ch}
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
