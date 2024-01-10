/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <key> <value>",
	Short: "Send AddItem command",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("both key and value are required")
		}

		if len(args) > 2 {
			return fmt.Errorf("required 2 arguments, got: %d", len(args))
		}

		client, err := NewClient()
		cobra.CheckErr(err)

		return client.AddItem(cmd.Context(), args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
