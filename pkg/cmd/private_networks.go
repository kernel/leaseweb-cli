package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var privateNetworksCmd = cli.Command{
	Name:    "private-networks",
	Aliases: []string{"pn"},
	Usage:   "Manage private networks",
	Commands: []*cli.Command{
		&pnListCmd,
		&pnGetCmd,
		&pnCreateCmd,
		&pnUpdateCmd,
		&pnDeleteCmd,
		&pnServersCmd,
	},
	HideHelpCommand: true,
}

var pnListCmd = cli.Command{
	Name:            "list",
	Usage:           "List private networks",
	Flags:           PaginationFlags,
	Action:          handlePNList,
	HideHelpCommand: true,
}

func handlePNList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/privateNetworks?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	networks := res.Get("privateNetworks")
	if !networks.Exists() || len(networks.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No private networks found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "NAME", "STATUS", "SUBNET", "LOCATION")
	networks.ForEach(func(_, n gjson.Result) bool {
		table.AddRow(
			n.Get("id").String(),
			n.Get("name").String(),
			n.Get("status").String(),
			n.Get("subnet").String(),
			n.Get("location.site").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var pnGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Get private network details",
	ArgsUsage:       "<network-id>",
	Action:          handlePNGet,
	HideHelpCommand: true,
}

func handlePNGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/privateNetworks/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var pnCreateCmd = cli.Command{
	Name:  "create",
	Usage: "Create a private network",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "Network name", Required: true},
	},
	Action:          handlePNCreate,
	HideHelpCommand: true,
}

func handlePNCreate(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"name": cmd.String("name")})
	res, err := client.PostJSON(ctx, "/bareMetals/v2/privateNetworks", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var pnUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Update a private network",
	ArgsUsage: "<network-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "New network name", Required: true},
	},
	Action:          handlePNUpdate,
	HideHelpCommand: true,
}

func handlePNUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"name": cmd.String("name")})
	res, err := client.PutJSON(ctx, "/bareMetals/v2/privateNetworks/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var pnDeleteCmd = cli.Command{
	Name:            "delete",
	Usage:           "Delete a private network",
	ArgsUsage:       "<network-id>",
	Action:          handlePNDelete,
	HideHelpCommand: true,
}

func handlePNDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/bareMetals/v2/privateNetworks/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted private network %s\n", args[0])
	return nil
}

var pnServersCmd = cli.Command{
	Name:      "servers",
	Usage:     "List servers in a private network",
	ArgsUsage: "<network-id>",
	Action:    handlePNServers,
	HideHelpCommand: true,
}

func handlePNServers(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/privateNetworks/"+args[0]+"/servers")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
