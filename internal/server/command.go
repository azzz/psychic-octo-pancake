package server

import (
	"fmt"
	"github.com/azzz/pillow/internal/omap"
	"io"
)

type Command interface {
	Exec(msg Message) error
}

type AddItemCommand struct {
	store *omap.OMap[string, string]
}

func (c AddItemCommand) Exec(msg Message) error {
	c.store.Set(msg.Key, msg.Value)
	return nil
}

type RemoveItemCommand struct {
	store *omap.OMap[string, string]
}

func (c RemoveItemCommand) Exec(msg Message) error {
	c.store.Delete(msg.Key)
	return nil
}

type GetItemCommand struct {
	w     io.Writer
	store *omap.OMap[string, string]
}

func (c GetItemCommand) Exec(msg Message) error {
	value, ok := c.store.Get(msg.Key)
	if !ok {
		return ValueNotFoundErr
	}

	item := Item{key: msg.Key, value: value}
	_, err := c.w.Write(item.Serialize())
	if err != nil {
		return fmt.Errorf("write item: %w", err)
	}

	return nil
}

type GetAllCommand struct {
	w     io.Writer
	store *omap.OMap[string, string]
}

func (c GetAllCommand) Exec(_ Message) error {
	for _, p := range c.store.Pairs() {
		item := Item{key: p.Key, value: p.Value}
		_, err := c.w.Write(item.Serialize())
		if err != nil {
			return fmt.Errorf("write item: %w", err)
		}
	}

	return nil
}
