package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var loadBalancersCmd = cli.Command{
	Name:    "load-balancers",
	Aliases: []string{"lb"},
	Usage:   "Manage public cloud load balancers",
	Commands: []*cli.Command{
		&lbListCmd,
		&lbGetCmd,
		&lbCreateCmd,
		&lbUpdateCmd,
		&lbDeleteCmd,
		&lbStartCmd,
		&lbStopCmd,
		&lbRebootCmd,
		&lbListenersCmd,
		&lbListenerCreateCmd,
		&lbListenerGetCmd,
		&lbListenerUpdateCmd,
		&lbListenerDeleteCmd,
		&lbIPsCmd,
		&lbIPGetCmd,
		&lbIPUpdateCmd,
		&lbIPNullCmd,
		&lbIPUnnullCmd,
		&lbMetricsCmd,
		&lbMonitoringEnableCmd,
		&lbMonitoringStatusCmd,
	},
	HideHelpCommand: true,
}

var lbListCmd = cli.Command{
	Name:            "list",
	Usage:           "List load balancers",
	Flags:           PaginationFlags,
	Action:          handleLBList,
	HideHelpCommand: true,
}

func handleLBList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/loadBalancers?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	lbs := res.Get("loadBalancers")
	if !lbs.Exists() || len(lbs.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No load balancers found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "REFERENCE", "TYPE", "REGION", "STATE", "IP")
	table.TruncOrder = []int{0}
	lbs.ForEach(func(_, lb gjson.Result) bool {
		ip := ""
		lb.Get("ips").ForEach(func(_, ipObj gjson.Result) bool {
			if ipObj.Get("version").Int() == 4 {
				ip = ipObj.Get("ip").String()
				return false
			}
			return true
		})
		table.AddRow(
			lb.Get("id").String(),
			lb.Get("reference").String(),
			lb.Get("type.name").String(),
			lb.Get("region").String(),
			lb.Get("state").String(),
			ip,
		)
		return true
	})
	table.Render()
	return nil
}

var lbGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Get load balancer details",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBGet,
	HideHelpCommand: true,
}

func handleLBGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/loadBalancers/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbCreateCmd = cli.Command{
	Name:  "create",
	Usage: "Create a load balancer",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "region", Usage: "Region", Required: true},
		&cli.StringFlag{Name: "type", Usage: "Load balancer type", Required: true},
		&cli.StringFlag{Name: "contract-type", Usage: "HOURLY or MONTHLY", Required: true},
		&cli.StringFlag{Name: "reference", Usage: "Reference name"},
	},
	Action:          handleLBCreate,
	HideHelpCommand: true,
}

func handleLBCreate(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"region":       cmd.String("region"),
		"type":         cmd.String("type"),
		"contractType": cmd.String("contract-type"),
	}
	if ref := cmd.String("reference"); ref != "" {
		payload["reference"] = ref
	}
	body, _ := json.Marshal(payload)
	res, err := client.PostJSON(ctx, "/publicCloud/v1/loadBalancers", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Update a load balancer",
	ArgsUsage: "<lb-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "New type"},
		&cli.StringFlag{Name: "reference", Usage: "New reference"},
	},
	Action:          handleLBUpdate,
	HideHelpCommand: true,
}

func handleLBUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]any{}
	if t := cmd.String("type"); t != "" {
		payload["type"] = t
	}
	if ref := cmd.String("reference"); ref != "" {
		payload["reference"] = ref
	}
	body, _ := json.Marshal(payload)
	res, err := client.PutJSON(ctx, "/publicCloud/v1/loadBalancers/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbDeleteCmd = cli.Command{
	Name:            "delete",
	Usage:           "Delete a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBDelete,
	HideHelpCommand: true,
}

func handleLBDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/publicCloud/v1/loadBalancers/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted load balancer %s\n", args[0])
	return nil
}

var lbListenersCmd = cli.Command{
	Name:            "listeners",
	Usage:           "List listeners for a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBListeners,
	HideHelpCommand: true,
}

func handleLBListeners(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/listeners")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbStartCmd = cli.Command{
	Name:            "start",
	Usage:           "Start a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBStart,
	HideHelpCommand: true,
}

func handleLBStart(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/start", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Started load balancer %s\n", args[0])
	return nil
}

var lbStopCmd = cli.Command{
	Name:            "stop",
	Usage:           "Stop a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBStop,
	HideHelpCommand: true,
}

func handleLBStop(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/stop", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Stopped load balancer %s\n", args[0])
	return nil
}

var lbRebootCmd = cli.Command{
	Name:            "reboot",
	Usage:           "Reboot a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBReboot,
	HideHelpCommand: true,
}

func handleLBReboot(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/reboot", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Rebooted load balancer %s\n", args[0])
	return nil
}

var lbListenerCreateCmd = cli.Command{
	Name:      "listener-create",
	Usage:     "Create a listener",
	ArgsUsage: "<lb-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload for the listener", Required: true},
	},
	Action:          handleLBListenerCreate,
	HideHelpCommand: true,
}

func handleLBListenerCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/listeners", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbListenerGetCmd = cli.Command{
	Name:            "listener-get",
	Usage:           "Get listener details",
	ArgsUsage:       "<lb-id> <listener-id>",
	Action:          handleLBListenerGet,
	HideHelpCommand: true,
}

func handleLBListenerGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and listener ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/listeners/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbListenerUpdateCmd = cli.Command{
	Name:      "listener-update",
	Usage:     "Update a listener",
	ArgsUsage: "<lb-id> <listener-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload for the listener", Required: true},
	},
	Action:          handleLBListenerUpdate,
	HideHelpCommand: true,
}

func handleLBListenerUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and listener ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/listeners/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbListenerDeleteCmd = cli.Command{
	Name:            "listener-delete",
	Usage:           "Delete a listener",
	ArgsUsage:       "<lb-id> <listener-id>",
	Action:          handleLBListenerDelete,
	HideHelpCommand: true,
}

func handleLBListenerDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and listener ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/listeners/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted listener %s\n", args[1])
	return nil
}

var lbIPsCmd = cli.Command{
	Name:            "ips",
	Usage:           "List IPs for a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBIPs,
	HideHelpCommand: true,
}

func handleLBIPs(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/ips")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbIPGetCmd = cli.Command{
	Name:            "ip-get",
	Usage:           "Get IP details for a load balancer",
	ArgsUsage:       "<lb-id> <ip>",
	Action:          handleLBIPGet,
	HideHelpCommand: true,
}

func handleLBIPGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/ips/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbIPUpdateCmd = cli.Command{
	Name:      "ip-update",
	Usage:     "Update IP for a load balancer",
	ArgsUsage: "<lb-id> <ip>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reverse-lookup", Usage: "Reverse lookup hostname", Required: true},
	},
	Action:          handleLBIPUpdate,
	HideHelpCommand: true,
}

func handleLBIPUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/ips/%s", args[0], args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbIPNullCmd = cli.Command{
	Name:            "ip-null",
	Usage:           "Null route an IP on a load balancer",
	ArgsUsage:       "<lb-id> <ip>",
	Action:          handleLBIPNull,
	HideHelpCommand: true,
}

func handleLBIPNull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/ips/%s/null", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s on %s\n", args[1], args[0])
	return nil
}

var lbIPUnnullCmd = cli.Command{
	Name:            "ip-unnull",
	Usage:           "Remove null route from an IP on a load balancer",
	ArgsUsage:       "<lb-id> <ip>",
	Action:          handleLBIPUnnull,
	HideHelpCommand: true,
}

func handleLBIPUnnull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("load balancer ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/ips/%s/unnull", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s on %s\n", args[1], args[0])
	return nil
}

var lbMetricsCmd = cli.Command{
	Name:      "metrics",
	Usage:     "Get load balancer metrics",
	ArgsUsage: "<lb-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "Metric type (cpu, datatraffic, connections, connectionsPerSecond, dataTransferred, dataTransferredPerSecond, requests, requestsPerSecond, responseCodes, responseCodesPerSecond)", Required: true},
		&cli.StringFlag{Name: "from", Usage: "Start date (ISO 8601)", Required: true},
		&cli.StringFlag{Name: "to", Usage: "End date (ISO 8601)", Required: true},
		&cli.StringFlag{Name: "granularity", Usage: "Granularity", Value: "1h"},
	},
	Action:          handleLBMetrics,
	HideHelpCommand: true,
}

func handleLBMetrics(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{
		"from":        cmd.String("from"),
		"to":          cmd.String("to"),
		"granularity": cmd.String("granularity"),
	})
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/loadBalancers/%s/metrics/%s%s", args[0], cmd.String("type"), q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var lbMonitoringEnableCmd = cli.Command{
	Name:            "monitoring-enable",
	Usage:           "Enable monitoring for a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBMonitoringEnable,
	HideHelpCommand: true,
}

func handleLBMonitoringEnable(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/monitoring/enable", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Enabled monitoring for %s\n", args[0])
	return nil
}

var lbMonitoringStatusCmd = cli.Command{
	Name:            "monitoring-status",
	Usage:           "Get monitoring status for a load balancer",
	ArgsUsage:       "<lb-id>",
	Action:          handleLBMonitoringStatus,
	HideHelpCommand: true,
}

func handleLBMonitoringStatus(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("load balancer ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/loadBalancers/"+args[0]+"/monitoring/status")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
