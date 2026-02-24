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
		&dsIPUpdateCmd,
		&dsIPNullCmd,
		&dsIPUnnullCmd,
		&dsPowerOnCmd,
		&dsPowerOffCmd,
		&dsPowerCycleCmd,
		&dsPowerStatusCmd,
		&dsRescueCmd,
		&dsRescueImagesCmd,
		&dsInstallCmd,
		&dsIPMIResetCmd,
		&dsCredentialsListCmd,
		&dsCredentialsGetCmd,
		&dsCredentialsCreateCmd,
		&dsCredentialsUpdateCmd,
		&dsCredentialsDeleteCmd,
		&dsJobsListCmd,
		&dsJobGetCmd,
		&dsJobCancelCmd,
		&dsJobExpireCmd,
		&dsJobRetryCmd,
		&dsHardwareInfoCmd,
		&dsHardwareMonitoringCmd,
		&dsHardwareMonitoringAllCmd,
		&dsHardwareScanCmd,
		&dsMetricsBandwidthCmd,
		&dsMetricsDatatrafficCmd,
		&dsNetworkInterfacesCmd,
		&dsLeasesListCmd,
		&dsLeasesCreateCmd,
		&dsLeasesDeleteCmd,
		&dsNullRouteHistoryCmd,
		&dsNotifBandwidthListCmd,
		&dsNotifBandwidthGetCmd,
		&dsNotifBandwidthCreateCmd,
		&dsNotifBandwidthUpdateCmd,
		&dsNotifBandwidthDeleteCmd,
		&dsNotifDatatrafficListCmd,
		&dsNotifDatatrafficGetCmd,
		&dsNotifDatatrafficCreateCmd,
		&dsNotifDatatrafficUpdateCmd,
		&dsNotifDatatrafficDeleteCmd,
		&dsNotifDDoSGetCmd,
		&dsNotifDDoSUpdateCmd,
		&dsOSListCmd,
		&dsOSGetCmd,
		&dsOSControlPanelsCmd,
		&dsControlPanelsCmd,
		&dsPrivateNetworkAddCmd,
		&dsPrivateNetworkRemoveCmd,
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

var dsIPUpdateCmd = cli.Command{
	Name:      "ip-update",
	Usage:     "Update an IP for a dedicated server",
	ArgsUsage: "<server-id> <ip>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reverse-lookup", Usage: "Reverse lookup hostname", Required: true},
	},
	Action:          handleDSIPUpdate,
	HideHelpCommand: true,
}

func handleDSIPUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/ips/%s", args[0], args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsIPNullCmd = cli.Command{
	Name:      "ip-null",
	Usage:     "Null route an IP on a dedicated server",
	ArgsUsage: "<server-id> <ip>",
	Action:    handleDSIPNull,
	HideHelpCommand: true,
}

func handleDSIPNull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/ips/%s/null", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s on %s\n", args[1], args[0])
	return nil
}

var dsIPUnnullCmd = cli.Command{
	Name:      "ip-unnull",
	Usage:     "Remove null route from an IP on a dedicated server",
	ArgsUsage: "<server-id> <ip>",
	Action:    handleDSIPUnnull,
	HideHelpCommand: true,
}

func handleDSIPUnnull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/ips/%s/unnull", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s on %s\n", args[1], args[0])
	return nil
}

var dsRescueImagesCmd = cli.Command{
	Name:            "rescue-images",
	Usage:           "List available rescue images",
	Action:          handleDSRescueImages,
	HideHelpCommand: true,
}

func handleDSRescueImages(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/rescueImages")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsIPMIResetCmd = cli.Command{
	Name:      "ipmi-reset",
	Usage:     "Launch IPMI reset for a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSIPMIReset,
	HideHelpCommand: true,
}

func handleDSIPMIReset(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/ipmiReset", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "IPMI reset initiated for %s\n", args[0])
	return nil
}

var dsCredentialsUpdateCmd = cli.Command{
	Name:      "credential-update",
	Usage:     "Update server credentials",
	ArgsUsage: "<server-id> <type> <username>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "password", Usage: "New password", Required: true},
	},
	Action:          handleDSCredentialsUpdate,
	HideHelpCommand: true,
}

func handleDSCredentialsUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("server ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"password": cmd.String("password")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/credentials/%s/%s", args[0], args[1], args[2]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsCredentialsDeleteCmd = cli.Command{
	Name:            "credential-delete",
	Usage:           "Delete server credentials",
	ArgsUsage:       "<server-id> <type> <username>",
	Action:          handleDSCredentialsDelete,
	HideHelpCommand: true,
}

func handleDSCredentialsDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("server ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted credential %s/%s for %s\n", args[1], args[2], args[0])
	return nil
}

var dsJobCancelCmd = cli.Command{
	Name:      "job-cancel",
	Usage:     "Cancel active job for a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSJobCancel,
	HideHelpCommand: true,
}

func handleDSJobCancel(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/cancelActiveJob", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Cancelled active job for %s\n", args[0])
	return nil
}

var dsJobExpireCmd = cli.Command{
	Name:      "job-expire",
	Usage:     "Expire active job for a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSJobExpire,
	HideHelpCommand: true,
}

func handleDSJobExpire(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/expireActiveJob", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Expired active job for %s\n", args[0])
	return nil
}

var dsJobRetryCmd = cli.Command{
	Name:      "job-retry",
	Usage:     "Retry a job",
	ArgsUsage: "<server-id> <job-id>",
	Action:    handleDSJobRetry,
	HideHelpCommand: true,
}

func handleDSJobRetry(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and job ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/jobs/%s/retry", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Retrying job %s for %s\n", args[1], args[0])
	return nil
}

var dsHardwareMonitoringCmd = cli.Command{
	Name:            "hardware-monitoring",
	Usage:           "Show hardware monitoring data for a server",
	ArgsUsage:       "<server-id>",
	Action:          handleDSHardwareMonitoring,
	HideHelpCommand: true,
}

func handleDSHardwareMonitoring(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/servers/"+args[0]+"/hardwareMonitoring")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsHardwareMonitoringAllCmd = cli.Command{
	Name:            "hardware-monitoring-all",
	Usage:           "Show hardware monitoring data for all servers",
	Action:          handleDSHardwareMonitoringAll,
	HideHelpCommand: true,
}

func handleDSHardwareMonitoringAll(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/hardwareMonitoring")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsHardwareScanCmd = cli.Command{
	Name:      "hardware-scan",
	Usage:     "Launch hardware scan for a dedicated server",
	ArgsUsage: "<server-id>",
	Action:    handleDSHardwareScan,
	HideHelpCommand: true,
}

func handleDSHardwareScan(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/bareMetals/v2/servers/"+args[0]+"/hardwareScan", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Hardware scan initiated for %s\n", args[0])
	return nil
}

var dsLeasesListCmd = cli.Command{
	Name:            "leases",
	Usage:           "List DHCP reservations for a dedicated server",
	ArgsUsage:       "<server-id>",
	Action:          handleDSLeasesList,
	HideHelpCommand: true,
}

func handleDSLeasesList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/servers/"+args[0]+"/leases")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsLeasesCreateCmd = cli.Command{
	Name:      "lease-create",
	Usage:     "Create a DHCP reservation for a dedicated server",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload for the reservation", Required: true},
	},
	Action:          handleDSLeasesCreate,
	HideHelpCommand: true,
}

func handleDSLeasesCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/bareMetals/v2/servers/"+args[0]+"/leases", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Created DHCP reservation for %s\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsLeasesDeleteCmd = cli.Command{
	Name:            "lease-delete",
	Usage:           "Delete a DHCP reservation for a dedicated server",
	ArgsUsage:       "<server-id>",
	Action:          handleDSLeasesDelete,
	HideHelpCommand: true,
}

func handleDSLeasesDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/bareMetals/v2/servers/"+args[0]+"/leases")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted DHCP reservation for %s\n", args[0])
	return nil
}

var dsNullRouteHistoryCmd = cli.Command{
	Name:            "null-route-history",
	Usage:           "Show null route history for a dedicated server",
	ArgsUsage:       "<server-id>",
	Action:          handleDSNullRouteHistory,
	HideHelpCommand: true,
}

func handleDSNullRouteHistory(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/servers/"+args[0]+"/nullRouteHistory")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifBandwidthListCmd = cli.Command{
	Name:            "notif-bandwidth-list",
	Usage:           "List bandwidth notification settings",
	ArgsUsage:       "<server-id>",
	Action:          handleDSNotifBandwidthList,
	HideHelpCommand: true,
}

func handleDSNotifBandwidthList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/bandwidth", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifBandwidthGetCmd = cli.Command{
	Name:            "notif-bandwidth-get",
	Usage:           "Get a bandwidth notification setting",
	ArgsUsage:       "<server-id> <notification-id>",
	Action:          handleDSNotifBandwidthGet,
	HideHelpCommand: true,
}

func handleDSNotifBandwidthGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/bandwidth/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifBandwidthCreateCmd = cli.Command{
	Name:      "notif-bandwidth-create",
	Usage:     "Create a bandwidth notification setting",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDSNotifBandwidthCreate,
	HideHelpCommand: true,
}

func handleDSNotifBandwidthCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/bandwidth", args[0]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifBandwidthUpdateCmd = cli.Command{
	Name:      "notif-bandwidth-update",
	Usage:     "Update a bandwidth notification setting",
	ArgsUsage: "<server-id> <notification-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDSNotifBandwidthUpdate,
	HideHelpCommand: true,
}

func handleDSNotifBandwidthUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/bandwidth/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifBandwidthDeleteCmd = cli.Command{
	Name:            "notif-bandwidth-delete",
	Usage:           "Delete a bandwidth notification setting",
	ArgsUsage:       "<server-id> <notification-id>",
	Action:          handleDSNotifBandwidthDelete,
	HideHelpCommand: true,
}

func handleDSNotifBandwidthDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/bandwidth/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted bandwidth notification %s for %s\n", args[1], args[0])
	return nil
}

var dsNotifDatatrafficListCmd = cli.Command{
	Name:            "notif-datatraffic-list",
	Usage:           "List data traffic notification settings",
	ArgsUsage:       "<server-id>",
	Action:          handleDSNotifDatatrafficList,
	HideHelpCommand: true,
}

func handleDSNotifDatatrafficList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/datatraffic", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifDatatrafficGetCmd = cli.Command{
	Name:            "notif-datatraffic-get",
	Usage:           "Get a data traffic notification setting",
	ArgsUsage:       "<server-id> <notification-id>",
	Action:          handleDSNotifDatatrafficGet,
	HideHelpCommand: true,
}

func handleDSNotifDatatrafficGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/datatraffic/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifDatatrafficCreateCmd = cli.Command{
	Name:      "notif-datatraffic-create",
	Usage:     "Create a data traffic notification setting",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDSNotifDatatrafficCreate,
	HideHelpCommand: true,
}

func handleDSNotifDatatrafficCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/datatraffic", args[0]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifDatatrafficUpdateCmd = cli.Command{
	Name:      "notif-datatraffic-update",
	Usage:     "Update a data traffic notification setting",
	ArgsUsage: "<server-id> <notification-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDSNotifDatatrafficUpdate,
	HideHelpCommand: true,
}

func handleDSNotifDatatrafficUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/datatraffic/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifDatatrafficDeleteCmd = cli.Command{
	Name:            "notif-datatraffic-delete",
	Usage:           "Delete a data traffic notification setting",
	ArgsUsage:       "<server-id> <notification-id>",
	Action:          handleDSNotifDatatrafficDelete,
	HideHelpCommand: true,
}

func handleDSNotifDatatrafficDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/datatraffic/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted data traffic notification %s for %s\n", args[1], args[0])
	return nil
}

var dsNotifDDoSGetCmd = cli.Command{
	Name:            "notif-ddos-get",
	Usage:           "Inspect DDoS notification settings",
	ArgsUsage:       "<server-id>",
	Action:          handleDSNotifDDoSGet,
	HideHelpCommand: true,
}

func handleDSNotifDDoSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/ddos", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsNotifDDoSUpdateCmd = cli.Command{
	Name:      "notif-ddos-update",
	Usage:     "Update DDoS notification settings",
	ArgsUsage: "<server-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDSNotifDDoSUpdate,
	HideHelpCommand: true,
}

func handleDSNotifDDoSUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/notificationSettings/ddos", args[0]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsOSListCmd = cli.Command{
	Name:            "os-list",
	Usage:           "List available operating systems",
	Flags:           PaginationFlags,
	Action:          handleDSOSList,
	HideHelpCommand: true,
}

func handleDSOSList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/operatingSystems?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsOSGetCmd = cli.Command{
	Name:            "os-get",
	Usage:           "Show an operating system",
	ArgsUsage:       "<os-id>",
	Action:          handleDSOSGet,
	HideHelpCommand: true,
}

func handleDSOSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("OS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/operatingSystems/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsOSControlPanelsCmd = cli.Command{
	Name:            "os-control-panels",
	Usage:           "List control panels for an operating system",
	ArgsUsage:       "<os-id>",
	Action:          handleDSOSControlPanels,
	HideHelpCommand: true,
}

func handleDSOSControlPanels(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("OS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/operatingSystems/"+args[0]+"/controlPanels")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsControlPanelsCmd = cli.Command{
	Name:            "control-panels",
	Usage:           "List all control panels",
	Action:          handleDSControlPanels,
	HideHelpCommand: true,
}

func handleDSControlPanels(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/bareMetals/v2/controlPanels")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var dsPrivateNetworkAddCmd = cli.Command{
	Name:            "private-network-add",
	Usage:           "Add a server to a private network",
	ArgsUsage:       "<server-id> <private-network-id>",
	Action:          handleDSPrivateNetworkAdd,
	HideHelpCommand: true,
}

func handleDSPrivateNetworkAdd(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and private network ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/privateNetworks/%s", args[0], args[1]), []byte("{}"))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Added server %s to private network %s\n", args[0], args[1])
	if res.Raw != "" {
		return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
	}
	return nil
}

var dsPrivateNetworkRemoveCmd = cli.Command{
	Name:            "private-network-remove",
	Usage:           "Remove a server from a private network",
	ArgsUsage:       "<server-id> <private-network-id>",
	Action:          handleDSPrivateNetworkRemove,
	HideHelpCommand: true,
}

func handleDSPrivateNetworkRemove(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("server ID and private network ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/bareMetals/v2/servers/%s/privateNetworks/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed server %s from private network %s\n", args[0], args[1])
	return nil
}
