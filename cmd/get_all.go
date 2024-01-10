package cmd

import (
	"github.com/spf13/cobra"
)

var getAllCmd = &cobra.Command{
	Use:   "get-all <key>",
	Short: "Send GetAllItems command",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewClient()
		cobra.CheckErr(err)

		return client.GetAllItems(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(getAllCmd)
}
