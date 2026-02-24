package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var apiKeysCmd = cli.Command{
	Name:    "api-keys",
	Aliases: []string{"keys"},
	Usage:   "Manage API keys",
	Commands: []*cli.Command{
		&akListCmd, &akCreateCmd, &akGetCmd, &akUpdateCmd, &akDeleteCmd,
		&akValidateCmd, &akCapabilitiesCmd,
	},
	HideHelpCommand: true,
}

var akListCmd = cli.Command{Name: "list", Usage: "List API keys", Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/auth/v2/apiKeys")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var akCreateCmd = cli.Command{Name: "create", Usage: "Create API key", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/auth/v2/apiKeys", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var akGetCmd = cli.Command{Name: "get", Usage: "Get API key", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("API key ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/auth/v2/apiKeys/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var akUpdateCmd = cli.Command{Name: "update", Usage: "Update API key", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("API key ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/auth/v2/apiKeys/"+args[0], []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var akDeleteCmd = cli.Command{Name: "delete", Usage: "Delete API key", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("API key ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/auth/v2/apiKeys/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted API key %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var akValidateCmd = cli.Command{Name: "validate", Usage: "Validate an API key", Flags: []cli.Flag{&cli.StringFlag{Name: "key", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"apiKey": cmd.String("key")})
	res, err := client.PostJSON(ctx, "/auth/v2/apiKeys/validate", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var akCapabilitiesCmd = cli.Command{Name: "capabilities", Usage: "List available capabilities", Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/auth/v2/apiKeys/capabilities")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
