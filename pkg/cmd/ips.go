package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var ipsCmd = cli.Command{
	Name:  "ips",
	Usage: "Manage IP addresses",
	Commands: []*cli.Command{
		&ipsListCmd,
		&ipsGetCmd,
		&ipsUpdateCmd,
		&ipsNullRouteCmd,
		&ipsRemoveNullRouteCmd,
		&ipsNullRouteHistoryCmd,
		&ipsNullRouteGetCmd,
		&ipsNullRouteUpdateCmd,
		&ipsNullRoutedIPv6Cmd,
		&ipsReverseLookupCmd,
		&ipsReverseLookupUpdateCmd,
	},
	HideHelpCommand: true,
}

var ipsListCmd = cli.Command{
	Name:  "list",
	Usage: "List IP addresses",
	Flags: append(PaginationFlags,
		&cli.StringFlag{Name: "version", Usage: "Filter by IP version (4 or 6)"},
		&cli.StringFlag{Name: "type", Usage: "Filter by type"},
		&cli.StringFlag{Name: "null-routed", Usage: "Filter by null routed status (true/false)"},
	),
	Action:          handleIPsList,
	HideHelpCommand: true,
}

func handleIPsList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	q := PaginationQuery(cmd)
	if v := cmd.String("version"); v != "" {
		q += "&version=" + v
	}
	if t := cmd.String("type"); t != "" {
		q += "&type=" + t
	}
	if nr := cmd.String("null-routed"); nr != "" {
		q += "&nullRouted=" + nr
	}

	res, err := client.Get(ctx, "/ipMgmt/v2/ips?"+q)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	ips := res.Get("ips")
	if !ips.Exists() || len(ips.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No IPs found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "IP", "VERSION", "TYPE", "REVERSE LOOKUP", "NULL ROUTED", "EQUIPMENT")
	table.TruncOrder = []int{3, 5}
	ips.ForEach(func(_, ip gjson.Result) bool {
		table.AddRow(
			ip.Get("ip").String(),
			fmt.Sprintf("v%d", ip.Get("version").Int()),
			ip.Get("type").String(),
			ip.Get("reverseLookup").String(),
			fmt.Sprintf("%t", ip.Get("nullRouted").Bool()),
			ip.Get("equipmentId").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var ipsGetCmd = cli.Command{
	Name:      "get",
	Usage:     "Get IP details",
	ArgsUsage: "<ip>",
	Action:    handleIPsGet,
	HideHelpCommand: true,
}

func handleIPsGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/ipMgmt/v2/ips/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Update IP (set reverse lookup)",
	ArgsUsage: "<ip>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reverse-lookup", Usage: "Reverse lookup hostname", Required: true},
	},
	Action:          handleIPsUpdate,
	HideHelpCommand: true,
}

func handleIPsUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{
		"reverseLookup": cmd.String("reverse-lookup"),
	})
	res, err := client.PutJSON(ctx, "/ipMgmt/v2/ips/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsNullRouteCmd = cli.Command{
	Name:      "null-route",
	Usage:     "Null route an IP",
	ArgsUsage: "<ip>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "comment", Usage: "Comment for the null route"},
	},
	Action:          handleIPsNullRoute,
	HideHelpCommand: true,
}

func handleIPsNullRoute(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	payload := map[string]string{}
	if c := cmd.String("comment"); c != "" {
		payload["comment"] = c
	}
	body, _ := json.Marshal(payload)
	res, err := client.PostJSON(ctx, "/ipMgmt/v2/ips/"+args[0]+"/nullRoute", body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s\n", args[0])
	if res.Raw != "" {
		return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
	}
	return nil
}

var ipsRemoveNullRouteCmd = cli.Command{
	Name:      "remove-null-route",
	Usage:     "Remove null route from an IP",
	ArgsUsage: "<ip>",
	Action:    handleIPsRemoveNullRoute,
	HideHelpCommand: true,
}

func handleIPsRemoveNullRoute(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/ipMgmt/v2/ips/"+args[0]+"/nullRoute")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s\n", args[0])
	return nil
}

var ipsNullRouteHistoryCmd = cli.Command{
	Name:  "null-route-history",
	Usage: "List null route history",
	Flags: PaginationFlags,
	Action:          handleIPsNullRouteHistory,
	HideHelpCommand: true,
}

func handleIPsNullRouteHistory(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/ipMgmt/v2/nullRoutes?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsNullRouteGetCmd = cli.Command{
	Name:            "null-route-get",
	Usage:           "Get null route details",
	ArgsUsage:       "<null-route-id>",
	Action:          handleIPsNullRouteGet,
	HideHelpCommand: true,
}

func handleIPsNullRouteGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("null route ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/ipMgmt/v2/nullRoutes/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsNullRouteUpdateCmd = cli.Command{
	Name:      "null-route-update",
	Usage:     "Update a null route",
	ArgsUsage: "<null-route-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "comment", Usage: "Comment for the null route"},
	},
	Action:          handleIPsNullRouteUpdate,
	HideHelpCommand: true,
}

func handleIPsNullRouteUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("null route ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"comment": cmd.String("comment")})
	res, err := client.PutJSON(ctx, "/ipMgmt/v2/nullRoutes/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsNullRoutedIPv6Cmd = cli.Command{
	Name:            "null-routed-ipv6",
	Usage:           "List null routed IPv6 addresses",
	ArgsUsage:       "<ip>",
	Action:          handleIPsNullRoutedIPv6,
	HideHelpCommand: true,
}

func handleIPsNullRoutedIPv6(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/ipMgmt/v2/ips/"+args[0]+"/nullRouted")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsReverseLookupCmd = cli.Command{
	Name:      "reverse-lookup",
	Usage:     "List reverse lookup records for an IPv6 range",
	ArgsUsage: "<ip>",
	Action:    handleIPsReverseLookup,
	HideHelpCommand: true,
}

func handleIPsReverseLookup(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/ipMgmt/v2/ips/"+args[0]+"/reverseLookup")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var ipsReverseLookupUpdateCmd = cli.Command{
	Name:      "reverse-lookup-update",
	Usage:     "Set or remove reverse lookup records for an IPv6 range",
	ArgsUsage: "<ip>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "records", Usage: "JSON array of reverse lookup records", Required: true},
	},
	Action:          handleIPsReverseLookupUpdate,
	HideHelpCommand: true,
}

func handleIPsReverseLookupUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("IP address required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/ipMgmt/v2/ips/"+args[0]+"/reverseLookup", []byte(cmd.String("records")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
