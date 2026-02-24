package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var trafficPolicyCmd = cli.Command{
	Name:    "traffic-policy",
	Aliases: []string{"tp"},
	Usage:   "Manage traffic policies",
	Commands: []*cli.Command{
		&tpListCmd, &tpGetCmd, &tpUpdateCmd,
		&tpHistoryCmd, &tpResetCmd,
	},
	HideHelpCommand: true,
}

var tpListCmd = cli.Command{Name: "list", Usage: "List traffic policies", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/trafficPolicy/v1/policies?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var tpGetCmd = cli.Command{Name: "get", Usage: "Get traffic policy", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("policy ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/trafficPolicy/v1/policies/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var tpUpdateCmd = cli.Command{Name: "update", Usage: "Update traffic policy", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("policy ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PatchJSON(ctx, "/trafficPolicy/v1/policies/"+args[0], []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var tpHistoryCmd = cli.Command{Name: "history", Usage: "List policy history", ArgsUsage: "<id>", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("policy ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/trafficPolicy/v1/policies/"+args[0]+"/history?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var tpResetCmd = cli.Command{Name: "reset", Usage: "Reset traffic policy counters", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("policy ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/trafficPolicy/v1/policies/"+args[0]+"/reset", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Reset policy %s\n", args[0])
	return nil
}, HideHelpCommand: true}
