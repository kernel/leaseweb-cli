package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var dedicatedRacksCmd = cli.Command{
	Name:    "dedicated-racks",
	Aliases: []string{"dr"},
	Usage:   "Manage dedicated racks",
	Commands: []*cli.Command{
		&drListCmd, &drGetCmd, &drUpdateCmd,
		&drCredentialsCmd, &drCredentialCreateCmd, &drCredentialsByTypeCmd,
		&drCredentialGetCmd, &drCredentialUpdateCmd, &drCredentialDeleteCmd,
		&drIPsCmd, &drIPGetCmd, &drIPUpdateCmd, &drIPNullCmd, &drIPUnnullCmd,
		&drMetricsBandwidthCmd, &drMetricsDatatrafficCmd,
		&drNullRouteHistoryCmd,
		&drNotifBandwidthListCmd, &drNotifBandwidthCreateCmd, &drNotifBandwidthGetCmd, &drNotifBandwidthUpdateCmd, &drNotifBandwidthDeleteCmd,
		&drNotifDatatrafficListCmd, &drNotifDatatrafficCreateCmd, &drNotifDatatrafficGetCmd, &drNotifDatatrafficUpdateCmd, &drNotifDatatrafficDeleteCmd,
		&drNotifDDoSGetCmd, &drNotifDDoSUpdateCmd,
	},
	HideHelpCommand: true,
}

func drPath(args []string) string { return "/bareMetals/v2/privateRacks/" + args[0] }

var drListCmd = cli.Command{Name: "list", Usage: "List dedicated racks", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/privateRacks?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drGetCmd = cli.Command{Name: "get", Usage: "Get dedicated rack", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, drPath(args))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drUpdateCmd = cli.Command{Name: "update", Usage: "Update", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "reference", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reference": cmd.String("reference")})
	res, err := client.PutJSON(ctx, drPath(args), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drCredentialsCmd = cli.Command{Name: "credentials", Usage: "List credentials", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, drPath(args)+"/credentials")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drCredentialCreateCmd = cli.Command{Name: "credential-create", Usage: "Create credential", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "type", Required: true}, &cli.StringFlag{Name: "username", Required: true}, &cli.StringFlag{Name: "password", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"type": cmd.String("type"), "username": cmd.String("username"), "password": cmd.String("password")})
	res, err := client.PostJSON(ctx, drPath(args)+"/credentials", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drCredentialsByTypeCmd = cli.Command{Name: "credentials-by-type", Usage: "List by type", ArgsUsage: "<id> <type>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("rack ID and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/credentials/%s", drPath(args), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drCredentialGetCmd = cli.Command{Name: "credential-get", Usage: "Get credential", ArgsUsage: "<id> <type> <username>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("rack ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/credentials/%s/%s", drPath(args), args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drCredentialUpdateCmd = cli.Command{Name: "credential-update", Usage: "Update credential", ArgsUsage: "<id> <type> <username>", Flags: []cli.Flag{&cli.StringFlag{Name: "password", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("rack ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"password": cmd.String("password")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("%s/credentials/%s/%s", drPath(args), args[1], args[2]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drCredentialDeleteCmd = cli.Command{Name: "credential-delete", Usage: "Delete credential", ArgsUsage: "<id> <type> <username>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("rack ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("%s/credentials/%s/%s", drPath(args), args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted credential %s/%s\n", args[1], args[2])
	return nil
}, HideHelpCommand: true}

var drIPsCmd = cli.Command{Name: "ips", Usage: "List IPs", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, drPath(args)+"/ips")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drIPGetCmd = cli.Command{Name: "ip-get", ArgsUsage: "<id> <ip>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("rack ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/ips/%s", drPath(args), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drIPUpdateCmd = cli.Command{Name: "ip-update", ArgsUsage: "<id> <ip>", Flags: []cli.Flag{&cli.StringFlag{Name: "reverse-lookup", Required: true}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("rack ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("%s/ips/%s", drPath(args), args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drIPNullCmd = cli.Command{Name: "ip-null", ArgsUsage: "<id> <ip>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("rack ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("%s/ips/%s/null", drPath(args), args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s\n", args[1])
	return nil
}, HideHelpCommand: true}

var drIPUnnullCmd = cli.Command{Name: "ip-unnull", ArgsUsage: "<id> <ip>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("rack ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("%s/ips/%s/unnull", drPath(args), args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s\n", args[1])
	return nil
}, HideHelpCommand: true}

var drMetricsBandwidthCmd = cli.Command{Name: "metrics-bandwidth", Usage: "Bandwidth metrics", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true}, &cli.StringFlag{Name: "aggregation", Value: "AVG"}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "aggregation": cmd.String("aggregation")})
	res, err := client.Get(ctx, fmt.Sprintf("%s/metrics/bandwidth%s", drPath(args), q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drMetricsDatatrafficCmd = cli.Command{Name: "metrics-datatraffic", Usage: "Datatraffic metrics", ArgsUsage: "<id>", Flags: []cli.Flag{&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true}, &cli.StringFlag{Name: "aggregation", Value: "SUM"}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "aggregation": cmd.String("aggregation")})
	res, err := client.Get(ctx, fmt.Sprintf("%s/metrics/datatraffic%s", drPath(args), q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drNullRouteHistoryCmd = cli.Command{Name: "null-route-history", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, drPath(args)+"/nullRouteHistory")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

func drNotifHandler(settingType, method string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		base := fmt.Sprintf("%s/notificationSettings/%s", drPath(args), settingType)
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		switch method {
		case "list":
			if len(args) < 1 {
				return fmt.Errorf("rack ID required")
			}
			res, err := client.Get(ctx, base)
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
		case "get":
			if len(args) < 2 {
				return fmt.Errorf("rack ID and notification ID required")
			}
			res, err := client.Get(ctx, base+"/"+args[1])
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
		case "create":
			if len(args) < 1 {
				return fmt.Errorf("rack ID required")
			}
			res, err := client.PostJSON(ctx, base, []byte(cmd.String("payload")))
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
		case "update":
			if len(args) < 2 {
				return fmt.Errorf("rack ID and notification ID required")
			}
			res, err := client.PutJSON(ctx, base+"/"+args[1], []byte(cmd.String("payload")))
			if err != nil {
				return err
			}
			return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
		case "delete":
			if len(args) < 2 {
				return fmt.Errorf("rack ID and notification ID required")
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

var payloadFlag = []cli.Flag{&cli.StringFlag{Name: "payload", Required: true}}

var drNotifBandwidthListCmd = cli.Command{Name: "notif-bandwidth-list", ArgsUsage: "<id>", Action: drNotifHandler("bandwidth", "list"), HideHelpCommand: true}
var drNotifBandwidthCreateCmd = cli.Command{Name: "notif-bandwidth-create", ArgsUsage: "<id>", Flags: payloadFlag, Action: drNotifHandler("bandwidth", "create"), HideHelpCommand: true}
var drNotifBandwidthGetCmd = cli.Command{Name: "notif-bandwidth-get", ArgsUsage: "<id> <nid>", Action: drNotifHandler("bandwidth", "get"), HideHelpCommand: true}
var drNotifBandwidthUpdateCmd = cli.Command{Name: "notif-bandwidth-update", ArgsUsage: "<id> <nid>", Flags: payloadFlag, Action: drNotifHandler("bandwidth", "update"), HideHelpCommand: true}
var drNotifBandwidthDeleteCmd = cli.Command{Name: "notif-bandwidth-delete", ArgsUsage: "<id> <nid>", Action: drNotifHandler("bandwidth", "delete"), HideHelpCommand: true}
var drNotifDatatrafficListCmd = cli.Command{Name: "notif-datatraffic-list", ArgsUsage: "<id>", Action: drNotifHandler("datatraffic", "list"), HideHelpCommand: true}
var drNotifDatatrafficCreateCmd = cli.Command{Name: "notif-datatraffic-create", ArgsUsage: "<id>", Flags: payloadFlag, Action: drNotifHandler("datatraffic", "create"), HideHelpCommand: true}
var drNotifDatatrafficGetCmd = cli.Command{Name: "notif-datatraffic-get", ArgsUsage: "<id> <nid>", Action: drNotifHandler("datatraffic", "get"), HideHelpCommand: true}
var drNotifDatatrafficUpdateCmd = cli.Command{Name: "notif-datatraffic-update", ArgsUsage: "<id> <nid>", Flags: payloadFlag, Action: drNotifHandler("datatraffic", "update"), HideHelpCommand: true}
var drNotifDatatrafficDeleteCmd = cli.Command{Name: "notif-datatraffic-delete", ArgsUsage: "<id> <nid>", Action: drNotifHandler("datatraffic", "delete"), HideHelpCommand: true}

var drNotifDDoSGetCmd = cli.Command{Name: "notif-ddos-get", ArgsUsage: "<id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, drPath(args)+"/notificationSettings/ddos")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var drNotifDDoSUpdateCmd = cli.Command{Name: "notif-ddos-update", ArgsUsage: "<id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("rack ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, drPath(args)+"/notificationSettings/ddos", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
