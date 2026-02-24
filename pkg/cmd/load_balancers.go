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
		&lbListenersCmd,
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
