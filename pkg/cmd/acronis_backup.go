package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var acronisBackupCmd = cli.Command{
	Name:    "acronis-backup",
	Aliases: []string{"backup"},
	Usage:   "Manage Acronis backup",
	Commands: []*cli.Command{
		&backupListCmd,
		&backupGetCmd,
		&backupMetricsCmd,
	},
	HideHelpCommand: true,
}

var backupListCmd = cli.Command{
	Name:            "list",
	Usage:           "List backup items",
	Action:          handleBackupList,
	HideHelpCommand: true,
}

func handleBackupList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/backup/v1/backup")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var backupGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Inspect a backup item",
	ArgsUsage:       "<equipment-id>",
	Action:          handleBackupGet,
	HideHelpCommand: true,
}

func handleBackupGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/backup/v1/backup/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var backupMetricsCmd = cli.Command{
	Name:            "metrics",
	Usage:           "Get storage usage metrics",
	Action:          handleBackupMetrics,
	HideHelpCommand: true,
}

func handleBackupMetrics(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/backup/v1/metrics/storage")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
