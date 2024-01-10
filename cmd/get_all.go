package cmd

import (
	"github.com/azzz/pillow/internal/client"
	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Use:   "get-all <key>",
	Short: "Send GetAllItems command",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New(config.AmqpUrl, config.AmqpQueue)
		cobra.CheckErr(err)

		return c.GetAllItems(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(getAllCmd)
}
