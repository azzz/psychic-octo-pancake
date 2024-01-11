package server

import (
	"bytes"
	"github.com/azzz/pillow/internal/omap"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newStore() *omap.OMap[string, string] {
	return omap.New[string, string]()
}

func TestAddItemCommand_Exec(t *testing.T) {
	store := newStore()
	cmd := AddItemCommand{store}
	assert.NoError(t, cmd.Exec(Message{Key: "Batman", Value: "DC"}))
	assert.NoError(t, cmd.Exec(Message{Key: "IronMan", Value: "Marvel"}))

	got1, _ := store.Get("Batman")
	assert.Equal(t, "DC", got1)

	got2, _ := store.Get("IronMan")
	assert.Equal(t, "Marvel", got2)
}

func TestRemoveItemCommand_Exec(t *testing.T) {
	store := newStore()
	store.Set("Batman", "DC")
	store.Set("SpiderMan", "DC")

	cmd := RemoveItemCommand{store}
	assert.NoError(t, cmd.Exec(Message{Key: "SpiderMan"}))

	got1, _ := store.Get("Batman")
	assert.Equal(t, "DC", got1)

	got2, ok2 := store.Get("SpiderMan")
	assert.False(t, ok2)
	assert.Equal(t, "", got2)
}

func TestGetItemCommand_Exec(t *testing.T) {
	store := newStore()
	store.Set("Batman", "DC")
	store.Set("SpiderMan", "Marvel")
	buf := bytes.NewBufferString("")

	cmd := GetItemCommand{buf, store}
	assert.NoError(t, cmd.Exec(Message{Key: "Batman"}))
	assert.Equalf(t, []string{"Batman", "SpiderMan"}, store.Keys(), "log is not expected to be changed")

	lines := []string{`"Batman":"DC"`, ""}
	assert.Equal(t, strings.Join(lines, "\n"), buf.String())
}

func TestGetAllCommand_Exec(t *testing.T) {
	store := newStore()
	store.Set("Batman", "DC")
	store.Set("SpiderMan", "Marvel")
	buf := bytes.NewBufferString("")

	cmd := GetAllCommand{buf, store}
	assert.NoError(t, cmd.Exec(Message{}))
	assert.Equalf(t, []string{"Batman", "SpiderMan"}, store.Keys(), "log is not expected to be changed")

	lines := []string{`"Batman":"DC"`, `"SpiderMan":"Marvel"`, ""}
	assert.Equal(t, strings.Join(lines, "\n"), buf.String())
}
