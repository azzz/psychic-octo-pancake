package cmd

import (
	"github.com/azzz/psychic-octo-pancake/internal/client"
	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Use:   "get-all <key>",
	Short: "Send GetAllItems command",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New(config.AmqpUrl, config.AmqpQueue)
		cobra.CheckErr(err)
		defer func() {
			logErr(c.Close(), "close client connections")
		}()

		return c.GetAllItems(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(getAllCmd)
}
