package cmd

import (
	"fmt"
	"github.com/azzz/pillow/internal/datalog"
	"github.com/azzz/pillow/internal/server"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run command server",
	RunE: func(cmd *cobra.Command, args []string) error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		defer func() {
			signal.Stop(sig)
		}()

		elog, err := datalog.Open(config.DataFile)
		if err != nil {
			return fmt.Errorf("open file storage: %w", err)
		}
		defer func() {
			logErr(elog.Close(), "close datalog")
		}()

		srv, err := server.New(config.AmqpUrl, config.AmqpQueue, elog)
		if err != nil {
			return fmt.Errorf("create server: %w", err)
		}

		go func() {
			select {
			case <-cmd.Context().Done():
				logger.Printf("ERROR: CLI dead by context")
			case <-sig:
				logger.Printf("INFO: handle CTL+C gracefully")
				if err := srv.Stop(); err != nil {
					logger.Printf("ERROR: stop server: %s", err)
				}
			}
		}()

		err = srv.Start(cmd.Context())
		return err
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
