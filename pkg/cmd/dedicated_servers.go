package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var dedicatedServersCmd = cli.Command{
	Name:    "dedicated-servers",
	Aliases: []string{"ds"},
	Usage:   "Manage dedicated servers",
	Commands: []*cli.Command{
		&dsListCmd,
		&dsGetCmd,
		&dsUpdateCmd,
		&dsIPsCmd,
		&dsIPGetCmd,
		&dsPowerOnCmd,
		&dsPowerOffCmd,
		&dsPowerCycleCmd,
		&dsPowerStatusCmd,
		&dsRescueCmd,
		&dsInstallCmd,
		&dsCredentialsListCmd,
		&dsCredentialsGetCmd,
		&dsCredentialsCreateCmd,
		&dsJobsListCmd,
		&dsJobGetCmd,
		&dsHardwareInfoCmd,
		&dsMetricsBandwidthCmd,
		&dsMetricsDatatrafficCmd,
		&dsNetworkInterfacesCmd,
	},
	HideHelpCommand: true,
}

var dsListCmd = cli.Command{
	Name:  "list",
	Usage: "List dedicated servers",
	Flags: append(PaginationFlags, &cli.StringFlag{
		Name:  "reference",
		Usage: "Filter by reference",
	}),
	Action:          handleDSList,
	HideHelpCommand: true,
}

func handleDSList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	q := PaginationQuery(cmd)
	if ref := cmd.String("reference"); ref != "" {
		q += "&reference=" + ref
	}

	res, err := client.Get(ctx, "/bareMetals/v2/servers?"+q)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	servers := res.Get("servers")
	if !servers.Exists() || !servers.IsArray() || len(servers.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No dedicated servers found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "REFERENCE", "SITE", "CHASSIS", "CPU", "RAM", "PUBLIC IP")
	table.TruncOrder = []int{3, 4, 6}
	servers.ForEach(func(_, s gjson.Result) bool {
		id := s.Get("id").String()
		ref := s.Get("reference").String()
		site := s.Get("location.site").String()
		chassis := s.Get("specs.chassis").String()
		cpu := s.Get("specs.cpu.type").String()
		ram := fmt.Sprintf("%d %s", s.Get("specs.ram.size").Int(), s.Get("specs.ram.unit").String())
		pubIP := s.Get("networkInterfaces.public.ip").String()
		table.AddRow(id, ref, site, chassis, cpu, ram, pubIP)
		return true
	})
	table.Render()
	return nil
}

var dsGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Get dedicated server details",
	ArgsUsage:       "<server-id>",
	Action:          handleDSGet,
	HideHelpCommand: true,
}

func handleDSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required\nUsage: lw dedicated-servers get <server-id>")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/servers/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Update dedicated server reference",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "reference",
			Usage:    "New reference string",
			Required: true,
		},
	},
	Action:          handleDSUpdate,
	HideHelpCommand: true,
}

func handleDSUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reference": cmd.String("reference")})
	res, err := client.PutJSON(ctx, "/bareMetals/v2/servers/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsIPsCmd = cli.Command{
	Name:      "ips",
	Usage:     "List IPs for a dedicated server",
	ArgsUsage: "<server-id>",
	Flags:     PaginationFlags,
	Action:    handleDSIPs,
	HideHelpCommand: true,
}

func handleDSIPs(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/ips?%s", args[0], PaginationQuery(cmd)))
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

	table := NewTableWriter(os.Stdout, "IP", "VERSION", "TYPE", "REVERSE LOOKUP", "NULL ROUTED")
	ips.ForEach(func(_, ip gjson.Result) bool {
		table.AddRow(
			ip.Get("ip").String(),
			fmt.Sprintf("v%d", ip.Get("version").Int()),
			ip.Get("type").String(),
			ip.Get("reverseLookup").String(),
			fmt.Sprintf("%t", ip.Get("nullRouted").Bool()),
		)
		return true
	})
	table.Render()
	return nil
}

var dsIPGetCmd = cli.Command{
	Name:      "ip-get",
	Usage:     "Get IP details for a dedicated server",
	ArgsUsage: "<server-id> <ip>",
	Action:    handleDSIPGet,
	HideHelpCommand: true,
}

func handleDSIPGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and IP required\nUsage: lw dedicated-servers ip-get <server-id> <ip>")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/ips/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsPowerOnCmd = cli.Command{
	Name:      "power-on",
	Usage:     "Power on a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSPowerOn,
	HideHelpCommand: true,
}

func handleDSPowerOn(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/powerOn", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Power on initiated for %s\n", args[0])
	return nil
}

var dsPowerOffCmd = cli.Command{
	Name:      "power-off",
	Usage:     "Power off a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSPowerOff,
	HideHelpCommand: true,
}

func handleDSPowerOff(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/powerOff", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Power off initiated for %s\n", args[0])
	return nil
}

var dsPowerCycleCmd = cli.Command{
	Name:      "power-cycle",
	Usage:     "Power cycle a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSPowerCycle,
	HideHelpCommand: true,
}

func handleDSPowerCycle(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/powerCycle", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Power cycle initiated for %s\n", args[0])
	return nil
}

var dsPowerStatusCmd = cli.Command{
	Name:      "power-status",
	Usage:     "Show power status of a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSPowerStatus,
	HideHelpCommand: true,
}

func handleDSPowerStatus(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/servers/"+args[0]+"/powerInfo")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsRescueCmd = cli.Command{
	Name:      "rescue",
	Usage:     "Launch rescue mode",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "os",
			Usage:    "Rescue image OS (e.g., RESCUE_GRML)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "power-cycle",
			Usage: "Power cycle after setting rescue (true/false)",
			Value: "true",
		},
	},
	Action:          handleDSRescue,
	HideHelpCommand: true,
}

func handleDSRescue(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]any{
		"rescueImageId": cmd.String("os"),
		"powerCycle":    cmd.String("power-cycle") == "true",
	})
	res, err := client.PostJSON(ctx, "/bareMetals/v2/servers/"+args[0]+"/rescueMode", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsInstallCmd = cli.Command{
	Name:      "install",
	Usage:     "Launch OS installation",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "os",
			Usage:    "Operating system ID (e.g., UBUNTU_22_04_64BIT)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "hostname",
			Usage: "Server hostname",
		},
	},
	Action:          handleDSInstall,
	HideHelpCommand: true,
}

func handleDSInstall(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]any{
		"operatingSystemId": cmd.String("os"),
	}
	if h := cmd.String("hostname"); h != "" {
		payload["hostname"] = h
	}
	body, _ := json.Marshal(payload)
	res, err := client.PostJSON(ctx, "/bareMetals/v2/servers/"+args[0]+"/install", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsCredentialsListCmd = cli.Command{
	Name:      "credentials",
	Usage:     "List credentials for a dedicated server",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "type",
			Usage: "Credential type filter (OPERATING_SYSTEM, RESCUE_MODE, etc.)",
		},
	},
	Action:          handleDSCredentialsList,
	HideHelpCommand: true,
}

func handleDSCredentialsList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/bareMetals/v2/servers/%s/credentials", args[0])
	if t := cmd.String("type"); t != "" {
		path += "/" + t
	}

	res, err := client.Get(ctx, path)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsCredentialsGetCmd = cli.Command{
	Name:      "credential-get",
	Usage:     "Get specific credentials",
	ArgsUsage: "<server-id> <type> <username>",
	Action:    handleDSCredentialsGet,
	HideHelpCommand: true,
}

func handleDSCredentialsGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("server ID, credential type, and username required\nUsage: lw dedicated-servers credential-get <server-id> <type> <username>")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsCredentialsCreateCmd = cli.Command{
	Name:      "credential-create",
	Usage:     "Create new credentials",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "Credential type", Required: true},
		&cli.StringFlag{Name: "username", Usage: "Username", Required: true},
		&cli.StringFlag{Name: "password", Usage: "Password", Required: true},
	},
	Action:          handleDSCredentialsCreate,
	HideHelpCommand: true,
}

func handleDSCredentialsCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{
		"type":     cmd.String("type"),
		"username": cmd.String("username"),
		"password": cmd.String("password"),
	})
	res, err := client.PostJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/credentials", args[0]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsJobsListCmd = cli.Command{
	Name:      "jobs",
	Usage:     "List jobs for a dedicated server",
	ArgsUsage: "<server-id>",
	Flags:     PaginationFlags,
	Action:    handleDSJobsList,
	HideHelpCommand: true,
}

func handleDSJobsList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/jobs?%s", args[0], PaginationQuery(cmd)))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	jobs := res.Get("jobs")
	if !jobs.Exists() || len(jobs.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No jobs found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "TYPE", "STATUS", "CREATED")
	jobs.ForEach(func(_, j gjson.Result) bool {
		table.AddRow(
			j.Get("uuid").String(),
			j.Get("type").String(),
			j.Get("status").String(),
			j.Get("createdAt").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var dsJobGetCmd = cli.Command{
	Name:      "job-get",
	Usage:     "Get job details",
	ArgsUsage: "<server-id> <job-id>",
	Action:    handleDSJobGet,
	HideHelpCommand: true,
}

func handleDSJobGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and job ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/jobs/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsHardwareInfoCmd = cli.Command{
	Name:      "hardware-info",
	Usage:     "Show hardware information",
	ArgsUsage: "<server-id>",
	Action:    handleDSHardwareInfo,
	HideHelpCommand: true,
}

func handleDSHardwareInfo(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/servers/"+args[0]+"/hardwareInfo")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsMetricsBandwidthCmd = cli.Command{
	Name:      "metrics-bandwidth",
	Usage:     "Show bandwidth metrics",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "from", Usage: "Start date (YYYY-MM-DD)", Required: true},
		&cli.StringFlag{Name: "to", Usage: "End date (YYYY-MM-DD)", Required: true},
		&cli.StringFlag{Name: "aggregation", Usage: "Aggregation (AVG, 95TH, SUM)", Value: "AVG"},
	},
	Action:          handleDSMetricsBandwidth,
	HideHelpCommand: true,
}

func handleDSMetricsBandwidth(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{
		"from":        cmd.String("from"),
		"to":          cmd.String("to"),
		"aggregation": cmd.String("aggregation"),
	})
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/metrics/bandwidth%s", args[0], q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsMetricsDatatrafficCmd = cli.Command{
	Name:      "metrics-datatraffic",
	Usage:     "Show datatraffic metrics",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "from", Usage: "Start date (YYYY-MM-DD)", Required: true},
		&cli.StringFlag{Name: "to", Usage: "End date (YYYY-MM-DD)", Required: true},
		&cli.StringFlag{Name: "aggregation", Usage: "Aggregation (SUM, AVG)", Value: "SUM"},
	},
	Action:          handleDSMetricsDatatraffic,
	HideHelpCommand: true,
}

func handleDSMetricsDatatraffic(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{
		"from":        cmd.String("from"),
		"to":          cmd.String("to"),
		"aggregation": cmd.String("aggregation"),
	})
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/metrics/datatraffic%s", args[0], q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNetworkInterfacesCmd = cli.Command{
	Name:      "network-interfaces",
	Usage:     "List network interfaces",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "action",
			Usage: "Perform action: open or close (on all interfaces)",
		},
		&cli.StringFlag{
			Name:  "interface",
			Usage: "Specific interface type (public, internal, remoteManagement)",
		},
	},
	Action:          handleDSNetworkInterfaces,
	HideHelpCommand: true,
}

func handleDSNetworkInterfaces(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	action := strings.ToLower(cmd.String("action"))
	iface := cmd.String("interface")

	if action != "" {
		path := fmt.Sprintf("/bareMetals/v2/servers/%s/networkInterfaces", args[0])
		if iface != "" {
			path += "/" + iface
		}
		path += "/" + action
		_, err := client.Post(ctx, path, "")
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Network interface %s action completed for %s\n", action, args[0])
		return nil
	}

	path := fmt.Sprintf("/bareMetals/v2/servers/%s/networkInterfaces", args[0])
	if iface != "" {
		path += "/" + iface
	}
	res, err := client.Get(ctx, path)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
