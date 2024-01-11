package server

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_buildCommand(t *testing.T) {
	type args struct {
		msg Message
	}
	tests := []struct {
		name    string
		args    args
		want    Command
		wantErr bool
	}{
		{
			"unsupported command",
			args{Message{
				Command: "makeMeBurger",
			}},
			nil,
			true,
		},

		{
			"AddItem",
			args{Message{Command: "AddItem"}},
			AddItemCommand{},
			false,
		},

		{
			"RemoveItem",
			args{Message{Command: "RemoveItem"}},
			RemoveItemCommand{},
			false,
		},

		{
			"GetAllItems",
			args{Message{Command: "GetAllItems"}},
			GetAllCommand{},
			false,
		},

		{
			"GetItem",
			args{Message{Command: "GetItem"}},
			GetItemCommand{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{}

			got, err := s.buildCommand(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_handle(t *testing.T) {
	t.Run("GetAllItems on empty log", func(t *testing.T) {
		var (
			buf = &bytes.Buffer{}
			s   = Server{
				store: newStore(),
				log:   buf,
			}
		)

		err := s.handle(message("GetAllItems"))
		assert.NoError(t, err)
		assert.Empty(t, buf)
	})

	t.Run("GetItem on empty log", func(t *testing.T) {
		var (
			buf = &bytes.Buffer{}
			s   = Server{
				store: newStore(),
				log:   buf,
			}
		)

		err := s.handle(message("GetItem", "Batman"))
		assert.Error(t, err)
		assert.Empty(t, buf)
	})

	t.Run("AddItem on empty KV store", func(t *testing.T) {
		var (
			buf = &bytes.Buffer{}
			s   = Server{
				store: newStore(),
				log:   buf,
			}
		)

		assert.NoError(t,
			s.handle(message("AddItem", "Batman", "DC")),
		)
		assert.NoError(t,
			s.handle(message("AddItem", "IronMan", "Marvel")),
		)

		assert.Equal(t, []string{"Batman", "IronMan"}, s.store.Keys())
		assert.Emptyf(t, buf, "should not write anything")
	})

	t.Run("GetItem", func(t *testing.T) {
		var (
			buf = &bytes.Buffer{}
			s   = Server{
				store: newStore(),
				log:   buf,
			}
		)

		s.store.Set("IronMan", "Marvel")
		s.store.Set("Batman", "DC")

		assert.NoError(t,
			s.handle(message("GetItem", "IronMan")),
		)

		lines := []string{`"IronMan":"Marvel"`, ""}
		assert.Equal(t, strings.Join(lines, "\n"), buf.String())
	})

	t.Run("GetAllItems", func(t *testing.T) {
		var (
			buf = &bytes.Buffer{}
			s   = Server{
				store: newStore(),
				log:   buf,
			}
		)

		s.store.Set("IronMan", "Marvel")
		s.store.Set("Batman", "DC")

		assert.NoError(t,
			s.handle(message("GetAllItems")),
		)

		lines := []string{`"IronMan":"Marvel"`, `"Batman":"DC"`, ""}
		assert.Equal(t, strings.Join(lines, "\n"), buf.String())
	})

	t.Run("RemoveItem", func(t *testing.T) {
		var (
			buf = &bytes.Buffer{}
			s   = Server{
				store: newStore(),
				log:   buf,
			}
		)

		s.store.Set("IronMan", "Marvel")
		s.store.Set("Batman", "DC")

		assert.NoError(t,
			s.handle(message("RemoveItem", "IronMan")),
		)

		assert.Emptyf(t, buf, "should not write anything")
		assert.Equal(t, []string{"Batman"}, s.store.Keys())
	})
}

func message(cmd string, args ...string) []byte {
	var (
		key   string
		value string
	)

	switch len(args) {
	case 0:
		break
	case 1:
		key = args[0]
	case 2:
		key = args[0]
		value = args[1]
	default:
		panic("too many arguments")
	}

	msg := Message{
		Command: cmd,
		Key:     key,
		Value:   value,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return data
}
