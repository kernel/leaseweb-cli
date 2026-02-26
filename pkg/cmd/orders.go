package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var ordersCmd = cli.Command{
	Name:  "orders",
	Usage: "Manage orders and product catalog",
	Commands: []*cli.Command{
		&ordersListCmd,
		&ordersGetCmd,
		&ordersProductsCmd,
	},
	HideHelpCommand: true,
}

// --- Account orders ---

var ordersListCmd = cli.Command{
	Name:            "list",
	Usage:           "List orders",
	Flags:           PaginationFlags,
	Action:          handleOrdersList,
	HideHelpCommand: true,
}

func handleOrdersList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/account/v1/orders?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("output")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	orders := res.Get("orders")
	if !orders.Exists() || len(orders.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No orders found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "TYPE", "ORIGIN", "CREATED AT")
	orders.ForEach(func(_, o gjson.Result) bool {
		table.AddRow(
			o.Get("id").String(),
			o.Get("type").String(),
			o.Get("origin").String(),
			dateOnly(o.Get("createdAt").String()),
		)
		return true
	})
	table.Render()
	return nil
}

var ordersGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Get order details",
	ArgsUsage:       "<order-id>",
	Action:          handleOrdersGet,
	HideHelpCommand: true,
}

func handleOrdersGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("order ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/account/v1/orders/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

// --- Product catalog & ordering ---

var ordersProductsCmd = cli.Command{
	Name:  "products",
	Usage: "Browse and order products",
	Commands: []*cli.Command{
		&ordersProductsDedicatedServersCmd,
		&ordersProductsVPSCmd,
	},
	HideHelpCommand: true,
}

// Dedicated servers

var ordersProductsDedicatedServersCmd = cli.Command{
	Name:    "dedicated-servers",
	Aliases: []string{"ds"},
	Usage:   "Browse and order dedicated servers",
	Commands: []*cli.Command{
		&ordersProductsDSListCmd,
		&ordersProductsDSGetCmd,
		&ordersProductsDSOrderCmd,
	},
	HideHelpCommand: true,
}

var ordersProductsDSListCmd = cli.Command{
	Name:  "list",
	Usage: "List available dedicated server configurations",
	Flags: append(PaginationFlags, []cli.Flag{
		&cli.StringFlag{Name: "location", Usage: "Filter by location"},
		&cli.StringFlag{Name: "ram", Usage: "Filter by RAM"},
		&cli.StringFlag{Name: "disk-size", Usage: "Filter by disk size"},
		&cli.StringFlag{Name: "disk-amount", Usage: "Filter by disk amount"},
	}...),
	Action:          handleProductsDSList,
	HideHelpCommand: true,
}

func handleProductsDSList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := PaginationQuery(cmd) + "&" + BuildQueryString(map[string]string{
		"location":   cmd.String("location"),
		"ram":        cmd.String("ram"),
		"diskSize":   cmd.String("disk-size"),
		"diskAmount": cmd.String("disk-amount"),
	})
	res, err := client.Get(ctx, "/ordering/v1/products/dedicatedServers?"+q)
	if err != nil {
		return err
	}

	format := cmd.Root().String("output")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	servers := res.Get("dedicatedServers")
	if !servers.Exists() || len(servers.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No dedicated server configurations found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "NAME", "CHASSIS", "CPU", "RAM", "STORAGE")
	table.TruncOrder = []int{1, 3, 5}
	servers.ForEach(func(_, s gjson.Result) bool {
		cpu := s.Get("cpu.quantity").String() + "x " + s.Get("cpu.speed").String()
		table.AddRow(
			s.Get("id").String(),
			s.Get("name").String(),
			s.Get("chassis").String(),
			cpu,
			s.Get("ram.amount").String()+" "+s.Get("ram.unit").String(),
			s.Get("storage.amount").String()+"x "+s.Get("storage.size").String()+" "+s.Get("storage.type").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var ordersProductsDSGetCmd = cli.Command{
	Name:      "get",
	Usage:     "Get dedicated server details and pricing",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "location", Usage: "Location for pricing", Required: true},
		&cli.BoolFlag{Name: "connected-to-aggregation-pool", Usage: "Include aggregation pool pricing"},
	},
	Action:          handleProductsDSGet,
	HideHelpCommand: true,
}

func handleProductsDSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{
		"location": cmd.String("location"),
	})
	if cmd.Bool("connected-to-aggregation-pool") {
		q += "&connectedToAggregationPool=true"
	}
	res, err := client.Get(ctx, "/ordering/v1/products/dedicatedServers/"+args[0]+q)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var ordersProductsDSOrderCmd = cli.Command{
	Name:      "order",
	Usage:     "Order a dedicated server",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "location", Usage: "Datacenter location (e.g. AMS-01)", Required: true},
		&cli.StringFlag{Name: "contract-term", Usage: "Contract term (1_MONTH, 3_MONTHS, 6_MONTHS, 1_YEAR, 2_YEARS, 3_YEARS)", Value: "1_MONTH"},
		&cli.BoolFlag{Name: "connected-to-aggregation-pool", Usage: "Connect to aggregation pool"},
	},
	Action:          handleProductsDSOrder,
	HideHelpCommand: true,
}

func handleProductsDSOrder(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"location":     cmd.String("location"),
		"contractTerm": cmd.String("contract-term"),
	}
	if cmd.Bool("connected-to-aggregation-pool") {
		payload["connectedToAggregationPool"] = true
	}
	body, _ := json.Marshal(payload)
	res, err := client.PostJSON(ctx, "/ordering/v1/products/dedicatedServers/"+args[0]+"/order", body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Order placed successfully. Order ID: %s\n", res.Get("orderId").String())
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

// VPS

var ordersProductsVPSCmd = cli.Command{
	Name:  "vps",
	Usage: "Browse and order VPS",
	Commands: []*cli.Command{
		&ordersProductsVPSListCmd,
		&ordersProductsVPSGetCmd,
		&ordersProductsVPSOrderCmd,
	},
	HideHelpCommand: true,
}

var ordersProductsVPSListCmd = cli.Command{
	Name:  "list",
	Usage: "List available VPS products",
	Flags: append(PaginationFlags, []cli.Flag{
		&cli.StringFlag{Name: "location", Usage: "Filter by location"},
	}...),
	Action:          handleProductsVPSList,
	HideHelpCommand: true,
}

func handleProductsVPSList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := PaginationQuery(cmd)
	if loc := cmd.String("location"); loc != "" {
		q += "&location=" + loc
	}
	res, err := client.Get(ctx, "/ordering/v1/products/vps?"+q)
	if err != nil {
		return err
	}

	format := cmd.Root().String("output")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	vpss := res.Get("vpss")
	if !vpss.Exists() || len(vpss.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No VPS products found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "NAME", "VCPU", "VRAM", "STORAGE", "TRAFFIC")
	vpss.ForEach(func(_, v gjson.Result) bool {
		table.AddRow(
			v.Get("id").String(),
			v.Get("name").String(),
			v.Get("vCpu").String(),
			v.Get("vRam").String(),
			v.Get("nvmeStorage").String(),
			v.Get("traffic").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var ordersProductsVPSGetCmd = cli.Command{
	Name:      "get",
	Usage:     "Get VPS product details and pricing",
	ArgsUsage: "<vps-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "location", Usage: "Location for pricing", Required: true},
		&cli.StringFlag{Name: "disk-upgrade", Usage: "Disk upgrade option"},
		&cli.StringFlag{Name: "operating-system", Usage: "Operating system option"},
		&cli.StringFlag{Name: "control-panel", Usage: "Control panel option"},
		&cli.StringFlag{Name: "contract-term", Usage: "Contract term"},
		&cli.StringFlag{Name: "billing-cycle", Usage: "Billing cycle"},
		&cli.StringFlag{Name: "sla", Usage: "Service level agreement"},
	},
	Action:          handleProductsVPSGet,
	HideHelpCommand: true,
}

func handleProductsVPSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{
		"location":              cmd.String("location"),
		"diskUpgrade":           cmd.String("disk-upgrade"),
		"operatingSystem":       cmd.String("operating-system"),
		"controlPanel":          cmd.String("control-panel"),
		"contractTerm":          cmd.String("contract-term"),
		"billingCycle":          cmd.String("billing-cycle"),
		"serviceLevelAgreement": cmd.String("sla"),
	})
	res, err := client.Get(ctx, "/ordering/v1/products/vps/"+args[0]+q)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var ordersProductsVPSOrderCmd = cli.Command{
	Name:      "order",
	Usage:     "Order a VPS",
	ArgsUsage: "<vps-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "location", Usage: "Datacenter location (e.g. AMS-01)", Required: true},
		&cli.StringFlag{Name: "disk-upgrade", Usage: "Disk upgrade option"},
		&cli.StringFlag{Name: "operating-system", Usage: "Operating system"},
		&cli.StringFlag{Name: "control-panel", Usage: "Control panel"},
		&cli.StringFlag{Name: "sla", Usage: "Service level agreement (Basic, Bronze, Silver, Gold, Platinum)"},
		&cli.StringFlag{Name: "contract-term", Usage: "Contract term (1_MONTH, 3_MONTHS, 6_MONTHS, 1_YEAR, 2_YEARS, 3_YEARS)", Value: "1_YEAR"},
		&cli.StringFlag{Name: "billing-cycle", Usage: "Billing cycle (1_MONTH, 3_MONTHS, 6_MONTHS, 1_YEAR)", Value: "1_MONTH"},
	},
	Action:          handleProductsVPSOrder,
	HideHelpCommand: true,
}

func handleProductsVPSOrder(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]string{
		"location":     cmd.String("location"),
		"contractTerm": cmd.String("contract-term"),
		"billingCycle": cmd.String("billing-cycle"),
	}
	if v := cmd.String("disk-upgrade"); v != "" {
		payload["diskUpgrade"] = v
	}
	if v := cmd.String("operating-system"); v != "" {
		payload["operatingSystem"] = v
	}
	if v := cmd.String("control-panel"); v != "" {
		payload["controlPanel"] = v
	}
	if v := cmd.String("sla"); v != "" {
		payload["serviceLevelAgreement"] = v
	}
	body, _ := json.Marshal(payload)
	res, err := client.PostJSON(ctx, "/ordering/v1/products/vps/"+args[0]+"/order", body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Order placed successfully. Order ID: %s\n", res.Get("orderId").String())
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
