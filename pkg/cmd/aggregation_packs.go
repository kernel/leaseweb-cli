package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var aggregationPacksCmd = cli.Command{
	Name:    "aggregation-packs",
	Aliases: []string{"ap"},
	Usage:   "Manage aggregation packs",
	Commands: []*cli.Command{
		&apListCmd,
		&apGetCmd,
	},
	HideHelpCommand: true,
}

var apListCmd = cli.Command{
	Name:            "list",
	Usage:           "List aggregation packs",
	Flags:           PaginationFlags,
	Action:          handleAPList,
	HideHelpCommand: true,
}

func handleAPList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/aggregationPacks?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var apGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Get aggregation pack details",
	ArgsUsage:       "<pack-id>",
	Action:          handleAPGet,
	HideHelpCommand: true,
}

func handleAPGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("aggregation pack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/aggregationPacks/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
