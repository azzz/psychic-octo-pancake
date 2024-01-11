package test

import (
	"context"
	"github.com/azzz/psychic-octo-pancake/internal/client"
	"github.com/azzz/psychic-octo-pancake/internal/server"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func amqpUrl() string {
	s := os.Getenv("AMQP_URL")
	if s == "" {
		log.Fatalf("AMQP_URL is not set")
	}
	return s
}

func amqpQueue() string {
	s := os.Getenv("AMQP_QUEUE")
	if s == "" {
		log.Fatalf("AMQP_QUEUE is not set")
	}
	return s
}

func setup(dataLog io.Writer) (*server.Server, *client.Client) {
	if err := godotenv.Load("../.env.test.env"); err != nil {
		log.Fatalf("load .env.test.env: %s", err)
	}

	srv, err := server.New(amqpUrl(), amqpQueue(), dataLog)
	if err != nil {
		log.Fatalf("create server: %s", err)
	}

	cl, err := client.New(amqpUrl(), amqpQueue())
	if err != nil {
		log.Fatalf("create client: %s", err)
	}

	return srv, cl
}

func Test(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "1" {
		t.Skip("Skip integrations tests. Use INTEGRATION_TESTS=1 env var to run them.")
	}

	dataLog := NewMemLog()
	ctx := context.Background()

	srv, client := setup(dataLog)
	defer func() { _ = srv.Stop() }()

	defer func() { _ = client.Close() }()

	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Fatalf("start server: %s", err)
		}
	}()

	require.NoError(t, client.AddItem(ctx, "france", "paris"))
	require.NoError(t, client.AddItem(ctx, "ukraine", "kyiv"))
	require.NoError(t, client.AddItem(ctx, "the netherlands", "hague"))

	require.NoError(t, client.GetItem(ctx, "france")) // print one key
	require.NoError(t, client.GetAllItems(ctx))       // print all keys
	require.NoError(t, client.AddItem(ctx, "the netherlands", "amsterdam"))
	require.NoError(t, client.RemoveItem(ctx, "france")) // remove key
	require.NoError(t, client.GetAllItems(ctx))          // print all: no france and updated netherlands

	// wait until server does all the things
	after := time.After(1 * time.Second)
	<-after

	lines := []string{
		// first GetItem
		`"france":"paris"`,
		// first GetAllItems
		`"france":"paris"`,
		`"ukraine":"kyiv"`,
		`"the netherlands":"hague"`,
		// second GetAllItems
		`"ukraine":"kyiv"`,
		`"the netherlands":"amsterdam"`,
		"",
	}

	assert.Equal(t, strings.Join(lines, "\n"), dataLog.String())
}
