package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var webhostingCmd = cli.Command{
	Name:    "webhosting",
	Aliases: []string{"wh"},
	Usage:   "Manage webhosting packages",
	Commands: []*cli.Command{
		&whListCmd, &whGetCmd, &whAvailableCmd,
		&whUsernamesCmd, &whUsernameGetCmd,
		&whDomainAliasesCmd, &whDomainAliasCreateCmd,
		&whCatchAllCmd, &whCatchAllUpdateCmd,
	},
	HideHelpCommand: true,
}

var whListCmd = cli.Command{Name: "list", Usage: "List webhosting packages", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/webhosting/v2/packages?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whGetCmd = cli.Command{Name: "get", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("package ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/webhosting/v2/packages/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whAvailableCmd = cli.Command{Name: "available", Usage: "List available webhosting packages", Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/webhosting/v2/packages/available")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whUsernamesCmd = cli.Command{Name: "usernames", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("package ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/webhosting/v2/packages/"+args[0]+"/usernames")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whUsernameGetCmd = cli.Command{Name: "username-get", ArgsUsage: "<id> <username>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("package ID and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/webhosting/v2/packages/%s/usernames/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whDomainAliasesCmd = cli.Command{Name: "domain-aliases", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("package ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/webhosting/v2/packages/"+args[0]+"/domainAliases")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whDomainAliasCreateCmd = cli.Command{Name: "domain-alias-create", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("package ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/webhosting/v2/packages/"+args[0]+"/domainAliases", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whCatchAllCmd = cli.Command{Name: "catch-all", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("package ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/webhosting/v2/packages/"+args[0]+"/catchAll")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var whCatchAllUpdateCmd = cli.Command{Name: "catch-all-update", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("package ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/webhosting/v2/packages/"+args[0]+"/catchAll", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
