package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var datacenterAccessCmd = cli.Command{
	Name:    "datacenter-access",
	Aliases: []string{"dca"},
	Usage:   "Manage datacenter access requests",
	Commands: []*cli.Command{
		&dcaListDatacentersCmd, &dcaListContactsCmd,
		&dcaListAccessCmd, &dcaCreateAccessCmd,
		&dcaGetAccessCmd, &dcaUpdateAccessCmd,
		&dcaDeleteAccessCmd, &dcaListVisitorsCmd,
	},
	HideHelpCommand: true,
}

var dcaListDatacentersCmd = cli.Command{Name: "datacenters", Usage: "List available datacenters", Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/datacenterAccess/v1/datacenters")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var dcaListContactsCmd = cli.Command{Name: "contacts", Usage: "List contacts", Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/datacenterAccess/v1/contacts")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var dcaListAccessCmd = cli.Command{Name: "list", Usage: "List access requests", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/datacenterAccess/v1/accessRequests?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var dcaCreateAccessCmd = cli.Command{Name: "create", Usage: "Create access request", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/datacenterAccess/v1/accessRequests", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var dcaGetAccessCmd = cli.Command{Name: "get", Usage: "Get access request", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("access request ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/datacenterAccess/v1/accessRequests/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var dcaUpdateAccessCmd = cli.Command{Name: "update", Usage: "Update access request", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("access request ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/datacenterAccess/v1/accessRequests/"+args[0], []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var dcaDeleteAccessCmd = cli.Command{Name: "delete", Usage: "Delete access request", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("access request ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/datacenterAccess/v1/accessRequests/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted access request %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var dcaListVisitorsCmd = cli.Command{Name: "visitors", Usage: "List visitors for request", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("access request ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/datacenterAccess/v1/accessRequests/"+args[0]+"/visitors")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
