package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var networkEquipmentCmd = cli.Command{
	Name:    "network-equipment",
	Aliases: []string{"ne"},
	Usage:   "Manage dedicated network equipment",
	Commands: []*cli.Command{
		&neListCmd, &neGetCmd, &neUpdateCmd,
		&neCredentialsCmd, &neCredentialCreateCmd, &neCredentialsByTypeCmd,
		&neCredentialGetCmd, &neCredentialUpdateCmd, &neCredentialDeleteCmd,
		&neIPsCmd, &neIPGetCmd, &neIPUpdateCmd, &neIPNullCmd, &neIPUnnullCmd,
		&neNullRouteHistoryCmd,
		&nePowerCycleCmd, &nePowerStatusCmd, &nePowerOffCmd, &nePowerOnCmd,
	},
	HideHelpCommand: true,
}

var neListCmd = cli.Command{Name: "list", Usage: "List network equipment", Flags: PaginationFlags, Action: handleNEList, HideHelpCommand: true}

func handleNEList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/networkEquipments?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neGetCmd = cli.Command{Name: "get", Usage: "Get network equipment", ArgsUsage: "<id>", Action: handleNEGet, HideHelpCommand: true}

func handleNEGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/networkEquipments/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neUpdateCmd = cli.Command{
	Name: "update", Usage: "Update network equipment", ArgsUsage: "<id>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "reference", Required: true}},
	Action: handleNEUpdate, HideHelpCommand: true,
}

func handleNEUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reference": cmd.String("reference")})
	res, err := client.PutJSON(ctx, "/bareMetals/v2/networkEquipments/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neCredentialsCmd = cli.Command{Name: "credentials", Usage: "List credentials", ArgsUsage: "<id>", Action: handleNECredentials, HideHelpCommand: true}

func handleNECredentials(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/credentials")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neCredentialCreateCmd = cli.Command{
	Name: "credential-create", Usage: "Create credentials", ArgsUsage: "<id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Required: true}, &cli.StringFlag{Name: "username", Required: true}, &cli.StringFlag{Name: "password", Required: true},
	},
	Action: handleNECredentialCreate, HideHelpCommand: true,
}

func handleNECredentialCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"type": cmd.String("type"), "username": cmd.String("username"), "password": cmd.String("password")})
	res, err := client.PostJSON(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/credentials", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neCredentialsByTypeCmd = cli.Command{Name: "credentials-by-type", Usage: "List credentials by type", ArgsUsage: "<id> <type>", Action: handleNECredentialsByType, HideHelpCommand: true}

func handleNECredentialsByType(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("network equipment ID and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/credentials/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neCredentialGetCmd = cli.Command{Name: "credential-get", Usage: "Get credential", ArgsUsage: "<id> <type> <username>", Action: handleNECredentialGet, HideHelpCommand: true}

func handleNECredentialGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neCredentialUpdateCmd = cli.Command{
	Name: "credential-update", Usage: "Update credential", ArgsUsage: "<id> <type> <username>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "password", Required: true}},
	Action: handleNECredentialUpdate, HideHelpCommand: true,
}

func handleNECredentialUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"password": cmd.String("password")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/credentials/%s/%s", args[0], args[1], args[2]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neCredentialDeleteCmd = cli.Command{Name: "credential-delete", Usage: "Delete credential", ArgsUsage: "<id> <type> <username>", Action: handleNECredentialDelete, HideHelpCommand: true}

func handleNECredentialDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted credential %s/%s\n", args[1], args[2])
	return nil
}

var neIPsCmd = cli.Command{Name: "ips", Usage: "List IPs", ArgsUsage: "<id>", Action: handleNEIPs, HideHelpCommand: true}

func handleNEIPs(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/ips")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neIPGetCmd = cli.Command{Name: "ip-get", Usage: "Get IP", ArgsUsage: "<id> <ip>", Action: handleNEIPGet, HideHelpCommand: true}

func handleNEIPGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/ips/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neIPUpdateCmd = cli.Command{
	Name: "ip-update", Usage: "Update IP", ArgsUsage: "<id> <ip>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "reverse-lookup", Required: true}},
	Action: handleNEIPUpdate, HideHelpCommand: true,
}

func handleNEIPUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/ips/%s", args[0], args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var neIPNullCmd = cli.Command{Name: "ip-null", Usage: "Null route IP", ArgsUsage: "<id> <ip>", Action: handleNEIPNull, HideHelpCommand: true}

func handleNEIPNull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/ips/%s/null", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s\n", args[1])
	return nil
}

var neIPUnnullCmd = cli.Command{Name: "ip-unnull", Usage: "Remove null route", ArgsUsage: "<id> <ip>", Action: handleNEIPUnnull, HideHelpCommand: true}

func handleNEIPUnnull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/bareMetals/v2/networkEquipments/%s/ips/%s/unnull", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s\n", args[1])
	return nil
}

var neNullRouteHistoryCmd = cli.Command{Name: "null-route-history", Usage: "Show null route history", ArgsUsage: "<id>", Action: handleNENullRouteHistory, HideHelpCommand: true}

func handleNENullRouteHistory(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("network equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/nullRouteHistory")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var nePowerCycleCmd = cli.Command{Name: "power-cycle", Usage: "Power cycle", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/powerCycle", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Power cycle initiated for %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var nePowerStatusCmd = cli.Command{Name: "power-status", Usage: "Show power status", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/powerInfo")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var nePowerOffCmd = cli.Command{Name: "power-off", Usage: "Power off", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/powerOff", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Powered off %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var nePowerOnCmd = cli.Command{Name: "power-on", Usage: "Power on", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/networkEquipments/"+args[0]+"/powerOn", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Powered on %s\n", args[0])
	return nil
}, HideHelpCommand: true}
