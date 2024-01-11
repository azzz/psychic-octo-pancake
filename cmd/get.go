package cmd

import (
	"errors"
	"fmt"
	"github.com/azzz/psychic-octo-pancake/internal/client"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Send GetItem command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing key")
		}

		if len(args) > 1 {
			return fmt.Errorf("required 1 argument, got: %d", len(args))
		}

		c, err := client.New(config.AmqpUrl, config.AmqpQueue)
		cobra.CheckErr(err)
		defer func() {
			logErr(c.Close(), "close client connections")
		}()

		return c.GetItem(cmd.Context(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

}
