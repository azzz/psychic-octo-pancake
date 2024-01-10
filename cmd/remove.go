package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <key>",
	Short: "Send RemoveItem command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing key")
		}

		if len(args) > 1 {
			return fmt.Errorf("required 1 argument, got: %d", len(args))
		}

		client, err := NewClient()
		cobra.CheckErr(err)

		return client.RemoveItem(cmd.Context(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
