package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var abuseReportsCmd = cli.Command{
	Name:    "abuse-reports",
	Aliases: []string{"abuse"},
	Usage:   "Manage abuse reports",
	Commands: []*cli.Command{
		&abuseListCmd, &abuseGetCmd, &abuseResolveCmd,
		&abuseMessagesCmd, &abuseMessageCreateCmd,
		&abuseAttachmentsCmd, &abuseAttachmentGetCmd, &abuseResolutionOptionsCmd,
	},
	HideHelpCommand: true,
}

func abusePath(id string) string { return "/abuse/v1/reports/" + id }

var abuseListCmd = cli.Command{Name: "list", Usage: "List abuse reports", Flags: append(PaginationFlags, &cli.StringFlag{Name: "status"}, &cli.StringFlag{Name: "sort-by"}), Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := PaginationQuery(cmd)
	if s := cmd.String("status"); s != "" {
		q += "&status=" + s
	}
	if s := cmd.String("sort-by"); s != "" {
		q += "&sortBy=" + s
	}
	res, err := client.Get(ctx, "/abuse/v1/reports?"+q)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var abuseGetCmd = cli.Command{Name: "get", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("report ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, abusePath(args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var abuseResolveCmd = cli.Command{Name: "resolve", Usage: "Resolve abuse report", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringSliceFlag{Name: "resolution", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("report ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string][]string{"resolutions": cmd.StringSlice("resolution")})
	_, err = client.PostJSON(ctx, abusePath(args[0])+"/resolve", body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Resolved report %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var abuseMessagesCmd = cli.Command{Name: "messages", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("report ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, abusePath(args[0])+"/messages")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var abuseMessageCreateCmd = cli.Command{Name: "message-create", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "body", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("report ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"body": cmd.String("body")})
	res, err := client.PostJSON(ctx, abusePath(args[0])+"/messages", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var abuseAttachmentsCmd = cli.Command{Name: "attachments", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("report ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, abusePath(args[0])+"/attachments")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var abuseAttachmentGetCmd = cli.Command{Name: "attachment-get", ArgsUsage: "<id> <attachment-id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("report ID and attachment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/attachments/%s", abusePath(args[0]), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var abuseResolutionOptionsCmd = cli.Command{Name: "resolution-options", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("report ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, abusePath(args[0])+"/resolutions")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
