package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var remoteManagementCmd = cli.Command{
	Name:    "remote-management",
	Aliases: []string{"rm"},
	Usage:   "Manage OpenVPN remote management",
	Commands: []*cli.Command{
		&rmProfilesCmd,
		&rmProfileGetCmd,
		&rmChangeCredentialsCmd,
	},
	HideHelpCommand: true,
}

var rmProfilesCmd = cli.Command{
	Name:            "profiles",
	Usage:           "List OpenVPN profiles",
	Action:          handleRMProfiles,
	HideHelpCommand: true,
}

func handleRMProfiles(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/remoteManagement/profiles")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var rmProfileGetCmd = cli.Command{
	Name:            "profile-get",
	Usage:           "Get an OpenVPN profile (e.g., lsw-rmvpn-ams-01.ovpn)",
	ArgsUsage:       "<profile-filename>",
	Action:          handleRMProfileGet,
	HideHelpCommand: true,
}

func handleRMProfileGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("profile filename required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/remoteManagement/profiles/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var rmChangeCredentialsCmd = cli.Command{
	Name:  "change-credentials",
	Usage: "Change OpenVPN credentials",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleRMChangeCredentials,
	HideHelpCommand: true,
}

func handleRMChangeCredentials(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/bareMetals/v2/remoteManagement/changeCredentials", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
