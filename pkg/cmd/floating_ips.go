package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var floatingIPsCmd = cli.Command{
	Name:    "floating-ips",
	Aliases: []string{"fip"},
	Usage:   "Manage floating IPs",
	Commands: []*cli.Command{
		&fipListCmd, &fipCreateCmd, &fipGetCmd, &fipUpdateCmd, &fipDeleteCmd,
		&fipDefinitionsCmd, &fipAssignCmd, &fipUnassignCmd,
	},
	HideHelpCommand: true,
}

var fipListCmd = cli.Command{Name: "list", Usage: "List floating IPs", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/floatingIps/v2/ranges?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var fipCreateCmd = cli.Command{Name: "create", Usage: "Create floating IP range", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/floatingIps/v2/ranges", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var fipGetCmd = cli.Command{Name: "get", Usage: "Get floating IP range", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("range ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/floatingIps/v2/ranges/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var fipUpdateCmd = cli.Command{Name: "update", Usage: "Update floating IP range", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "comment"}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("range ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"comment": cmd.String("comment")})
	res, err := client.PutJSON(ctx, "/floatingIps/v2/ranges/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var fipDeleteCmd = cli.Command{Name: "delete", Usage: "Delete floating IP range", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("range ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/floatingIps/v2/ranges/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted range %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var fipDefinitionsCmd = cli.Command{Name: "definitions", Usage: "List floating IP definitions", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/floatingIps/v2/ranges/definitions?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var fipAssignCmd = cli.Command{Name: "assign", Usage: "Assign floating IP to server", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "anchor-ip", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("range ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"anchorIp": cmd.String("anchor-ip")})
	res, err := client.PostJSON(ctx, "/floatingIps/v2/ranges/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var fipUnassignCmd = cli.Command{Name: "unassign", Usage: "Remove floating IP assignment", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("range ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/floatingIps/v2/ranges/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Unassigned range %s\n", args[0])
	return nil
}, HideHelpCommand: true}
