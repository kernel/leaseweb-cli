package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var colocationsCmd = cli.Command{
	Name:    "colocations",
	Aliases: []string{"colo"},
	Usage:   "Manage colocations",
	Commands: []*cli.Command{
		&coloListCmd, &coloGetCmd, &coloUpdateCmd,
		&coloCredentialsCmd, &coloCredentialCreateCmd, &coloCredentialsByTypeCmd,
		&coloCredentialGetCmd, &coloCredentialUpdateCmd, &coloCredentialDeleteCmd,
		&coloIPsCmd, &coloIPGetCmd, &coloIPUpdateCmd, &coloIPNullCmd, &coloIPUnnullCmd,
		&coloMetricsBandwidthCmd, &coloMetricsDatatrafficCmd,
		&coloNetworkInterfaceCmd, &coloNetworkInterfaceActionCmd,
		&coloNullRouteHistoryCmd,
		&coloNotifBandwidthListCmd, &coloNotifBandwidthCreateCmd, &coloNotifBandwidthGetCmd, &coloNotifBandwidthUpdateCmd, &coloNotifBandwidthDeleteCmd,
		&coloNotifDatatrafficListCmd, &coloNotifDatatrafficCreateCmd, &coloNotifDatatrafficGetCmd, &coloNotifDatatrafficUpdateCmd, &coloNotifDatatrafficDeleteCmd,
		&coloNotifDDoSGetCmd, &coloNotifDDoSUpdateCmd,
	},
	HideHelpCommand: true,
}

func coloPath(args []string) string { return "/bareMetals/v2/colocations/" + args[0] }

func coloSimpleGet(subpath string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < 1 {
			return fmt.Errorf("colocation ID required")
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		res, err := client.Get(ctx, coloPath(args)+subpath)
		if err != nil {
			return err
		}
		return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
	}
}

var coloListCmd = cli.Command{Name: "list", Usage: "List colocations", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/colocations?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloGetCmd = cli.Command{Name: "get", ArgsUsage: "<id>", Action: coloSimpleGet(""), HideHelpCommand: true}

var coloUpdateCmd = cli.Command{Name: "update", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "reference", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("colocation ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reference": cmd.String("reference")})
	res, err := client.PutJSON(ctx, coloPath(args), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloCredentialsCmd = cli.Command{Name: "credentials", ArgsUsage: "<id>", Action: coloSimpleGet("/credentials"), HideHelpCommand: true}

var coloCredentialCreateCmd = cli.Command{Name: "credential-create", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "type", Required: true}, &cli.StringFlag{Name: "username", Required: true}, &cli.StringFlag{Name: "password", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("colocation ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"type": cmd.String("type"), "username": cmd.String("username"), "password": cmd.String("password")})
	res, err := client.PostJSON(ctx, coloPath(args)+"/credentials", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloCredentialsByTypeCmd = cli.Command{Name: "credentials-by-type", ArgsUsage: "<id> <type>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("colocation ID and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/credentials/%s", coloPath(args), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloCredentialGetCmd = cli.Command{Name: "credential-get", ArgsUsage: "<id> <type> <username>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("colocation ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/credentials/%s/%s", coloPath(args), args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloCredentialUpdateCmd = cli.Command{Name: "credential-update", ArgsUsage: "<id> <type> <username>", Flags: []cli.Flag{&cli.StringFlag{Name: "password", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("colocation ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"password": cmd.String("password")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("%s/credentials/%s/%s", coloPath(args), args[1], args[2]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloCredentialDeleteCmd = cli.Command{Name: "credential-delete", ArgsUsage: "<id> <type> <username>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("colocation ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("%s/credentials/%s/%s", coloPath(args), args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted credential %s/%s\n", args[1], args[2])
	return nil
}, HideHelpCommand: true}

var coloIPsCmd = cli.Command{Name: "ips", ArgsUsage: "<id>", Action: coloSimpleGet("/ips"), HideHelpCommand: true}

var coloIPGetCmd = cli.Command{Name: "ip-get", ArgsUsage: "<id> <ip>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("colocation ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/ips/%s", coloPath(args), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloIPUpdateCmd = cli.Command{Name: "ip-update", ArgsUsage: "<id> <ip>", Flags: []cli.Flag{&cli.StringFlag{Name: "reverse-lookup", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("colocation ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("%s/ips/%s", coloPath(args), args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloIPNullCmd = cli.Command{Name: "ip-null", ArgsUsage: "<id> <ip>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("colocation ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("%s/ips/%s/null", coloPath(args), args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s\n", args[1])
	return nil
}, HideHelpCommand: true}

var coloIPUnnullCmd = cli.Command{Name: "ip-unnull", ArgsUsage: "<id> <ip>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("colocation ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("%s/ips/%s/unnull", coloPath(args), args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s\n", args[1])
	return nil
}, HideHelpCommand: true}

var coloMetricsBandwidthCmd = cli.Command{Name: "metrics-bandwidth", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true}, &cli.StringFlag{Name: "aggregation", Value: "AVG"}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("colocation ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "aggregation": cmd.String("aggregation")})
	res, err := client.Get(ctx, fmt.Sprintf("%s/metrics/bandwidth%s", coloPath(args), q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloMetricsDatatrafficCmd = cli.Command{Name: "metrics-datatraffic", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true}, &cli.StringFlag{Name: "aggregation", Value: "SUM"}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("colocation ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "aggregation": cmd.String("aggregation")})
	res, err := client.Get(ctx, fmt.Sprintf("%s/metrics/datatraffic%s", coloPath(args), q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var coloNetworkInterfaceCmd = cli.Command{Name: "network-interface", Usage: "Inspect public network interface", ArgsUsage: "<id>", Action: coloSimpleGet("/networkInterfaces/public"), HideHelpCommand: true}

var coloNetworkInterfaceActionCmd = cli.Command{Name: "network-interface-action", Usage: "Open or close public network interface", ArgsUsage: "<id> <action>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("colocation ID and action (open/close) required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("%s/networkInterfaces/public/%s", coloPath(args), args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Network interface %s completed for %s\n", args[1], args[0])
	return nil
}, HideHelpCommand: true}

var coloNullRouteHistoryCmd = cli.Command{Name: "null-route-history", ArgsUsage: "<id>", Action: coloSimpleGet("/nullRouteHistory"), HideHelpCommand: true}

func coloNotifHandler(settingType, method string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		base := fmt.Sprintf("%s/notificationSettings/%s", coloPath(args), settingType)
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		switch method {
		case "list":
			if len(args) < 1 {
				return fmt.Errorf("colocation ID required")
			}
			res, err := client.Get(ctx, base)
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
		case "get":
			if len(args) < 2 {
				return fmt.Errorf("colocation ID and notification ID required")
			}
			res, err := client.Get(ctx, base+"/"+args[1])
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
		case "create":
			if len(args) < 1 {
				return fmt.Errorf("colocation ID required")
			}
			res, err := client.PostJSON(ctx, base, []byte(cmd.String("payload")))
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
		case "update":
			if len(args) < 2 {
				return fmt.Errorf("colocation ID and notification ID required")
			}
			res, err := client.PutJSON(ctx, base+"/"+args[1], []byte(cmd.String("payload")))
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
		case "delete":
			if len(args) < 2 {
				return fmt.Errorf("colocation ID and notification ID required")
			}
			_, err := client.Delete(ctx, base+"/"+args[1])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Deleted notification %s\n", args[1])
			return nil
		}
		return nil
	}
}

var coloNotifBandwidthListCmd = cli.Command{Name: "notif-bandwidth-list", ArgsUsage: "<id>", Action: coloNotifHandler("bandwidth", "list"), HideHelpCommand: true}
var coloNotifBandwidthCreateCmd = cli.Command{Name: "notif-bandwidth-create", ArgsUsage: "<id>", Flags: payloadFlag, Action: coloNotifHandler("bandwidth", "create"), HideHelpCommand: true}
var coloNotifBandwidthGetCmd = cli.Command{Name: "notif-bandwidth-get", ArgsUsage: "<id> <nid>", Action: coloNotifHandler("bandwidth", "get"), HideHelpCommand: true}
var coloNotifBandwidthUpdateCmd = cli.Command{Name: "notif-bandwidth-update", ArgsUsage: "<id> <nid>", Flags: payloadFlag, Action: coloNotifHandler("bandwidth", "update"), HideHelpCommand: true}
var coloNotifBandwidthDeleteCmd = cli.Command{Name: "notif-bandwidth-delete", ArgsUsage: "<id> <nid>", Action: coloNotifHandler("bandwidth", "delete"), HideHelpCommand: true}
var coloNotifDatatrafficListCmd = cli.Command{Name: "notif-datatraffic-list", ArgsUsage: "<id>", Action: coloNotifHandler("datatraffic", "list"), HideHelpCommand: true}
var coloNotifDatatrafficCreateCmd = cli.Command{Name: "notif-datatraffic-create", ArgsUsage: "<id>", Flags: payloadFlag, Action: coloNotifHandler("datatraffic", "create"), HideHelpCommand: true}
var coloNotifDatatrafficGetCmd = cli.Command{Name: "notif-datatraffic-get", ArgsUsage: "<id> <nid>", Action: coloNotifHandler("datatraffic", "get"), HideHelpCommand: true}
var coloNotifDatatrafficUpdateCmd = cli.Command{Name: "notif-datatraffic-update", ArgsUsage: "<id> <nid>", Flags: payloadFlag, Action: coloNotifHandler("datatraffic", "update"), HideHelpCommand: true}
var coloNotifDatatrafficDeleteCmd = cli.Command{Name: "notif-datatraffic-delete", ArgsUsage: "<id> <nid>", Action: coloNotifHandler("datatraffic", "delete"), HideHelpCommand: true}

var coloNotifDDoSGetCmd = cli.Command{Name: "notif-ddos-get", ArgsUsage: "<id>", Action: coloSimpleGet("/notificationSettings/ddos"), HideHelpCommand: true}
var coloNotifDDoSUpdateCmd = cli.Command{Name: "notif-ddos-update", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("colocation ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, coloPath(args)+"/notificationSettings/ddos", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
