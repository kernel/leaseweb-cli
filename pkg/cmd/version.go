package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var Version = "dev"

var versionCmd = cli.Command{
	Name:            "version",
	Usage:           "Print the version",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Fprintf(os.Stdout, "lw version %s\n", Version)
		return nil
	},
	HideHelpCommand: true,
}
