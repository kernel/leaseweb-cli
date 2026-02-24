package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var servicesCmd = cli.Command{
	Name:  "services",
	Usage: "Manage services",
	Commands: []*cli.Command{
		&servicesListCmd,
		&servicesGetCmd,
		&servicesUpdateCmd,
		&servicesCancelCmd,
		&servicesUncancelCmd,
		&servicesCancellationReasonsCmd,
	},
	HideHelpCommand: true,
}

var servicesListCmd = cli.Command{
	Name:  "list",
	Usage: "List services",
	Flags: PaginationFlags,
	Action:          handleServicesList,
	HideHelpCommand: true,
}

func handleServicesList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/services/v1/services?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	services := res.Get("services")
	if !services.Exists() || len(services.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No services found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "REFERENCE", "PRODUCT", "STATUS", "START DATE", "END DATE")
	table.TruncOrder = []int{2, 1}
	services.ForEach(func(_, svc gjson.Result) bool {
		table.AddRow(
			svc.Get("id").String(),
			svc.Get("reference").String(),
			svc.Get("productId").String(),
			svc.Get("status").String(),
			svc.Get("startDate").String(),
			svc.Get("endDate").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var servicesGetCmd = cli.Command{
	Name:      "get",
	Usage:     "Get service details",
	ArgsUsage: "<service-id>",
	Action:    handleServicesGet,
	HideHelpCommand: true,
}

func handleServicesGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("service ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/services/v1/services/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var servicesUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Update a service",
	ArgsUsage: "<service-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reference", Usage: "New reference", Required: true},
	},
	Action:          handleServicesUpdate,
	HideHelpCommand: true,
}

func handleServicesUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("service ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{
		"reference": cmd.String("reference"),
	})
	res, err := client.PutJSON(ctx, "/services/v1/services/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var servicesCancelCmd = cli.Command{
	Name:      "cancel",
	Usage:     "Cancel a service",
	ArgsUsage: "<service-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reason", Usage: "Cancellation reason", Required: true},
		&cli.StringFlag{Name: "reason-detail", Usage: "Cancellation reason detail"},
	},
	Action:          handleServicesCancel,
	HideHelpCommand: true,
}

func handleServicesCancel(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("service ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]string{
		"reason": cmd.String("reason"),
	}
	if d := cmd.String("reason-detail"); d != "" {
		payload["reasonDetail"] = d
	}
	body, _ := json.Marshal(payload)
	_, err = client.PostJSON(ctx, "/services/v1/services/"+args[0]+"/cancel", body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Cancelled service %s\n", args[0])
	return nil
}

var servicesUncancelCmd = cli.Command{
	Name:      "uncancel",
	Usage:     "Uncancel a service",
	ArgsUsage: "<service-id>",
	Action:    handleServicesUncancel,
	HideHelpCommand: true,
}

func handleServicesUncancel(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("service ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/services/v1/services/"+args[0]+"/uncancel", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Uncancelled service %s\n", args[0])
	return nil
}

var servicesCancellationReasonsCmd = cli.Command{
	Name:            "cancellation-reasons",
	Usage:           "List cancellation reasons",
	Action:          handleServicesCancellationReasons,
	HideHelpCommand: true,
}

func handleServicesCancellationReasons(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/services/v1/services/cancellationReasons")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
