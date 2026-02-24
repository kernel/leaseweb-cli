package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var instancesCmd = cli.Command{
	Name:    "instances",
	Aliases: []string{"i"},
	Usage:   "Manage public cloud instances",
	Commands: []*cli.Command{
		&instancesListCmd,
		&instancesGetCmd,
		&instancesLaunchCmd,
		&instancesTerminateCmd,
		&instancesStartCmd,
		&instancesStopCmd,
		&instancesRebootCmd,
		&instancesUpdateCmd,
		&instancesConsoleCmd,
		&instancesCredentialsCmd,
		&instancesIPsCmd,
		&instancesSnapshotsListCmd,
		&instancesSnapshotCreateCmd,
		&instancesMetricsCmd,
		&instancesRegionsCmd,
		&instancesTypesCmd,
		&instancesImagesCmd,
	},
	HideHelpCommand: true,
}

var instancesListCmd = cli.Command{
	Name:  "list",
	Usage: "List public cloud instances",
	Flags: PaginationFlags,
	Action:          handleInstancesList,
	HideHelpCommand: true,
}

func handleInstancesList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/instances?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	instances := res.Get("instances")
	if !instances.Exists() || len(instances.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No instances found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "ID", "REFERENCE", "TYPE", "REGION", "STATE", "IMAGE", "PUBLIC IP")
	table.TruncOrder = []int{0, 5, 6}
	instances.ForEach(func(_, inst gjson.Result) bool {
		id := inst.Get("id").String()
		ref := inst.Get("reference").String()
		itype := inst.Get("type.name").String()
		region := inst.Get("region").String()
		state := inst.Get("state").String()
		image := inst.Get("image.id").String()
		ip := ""
		ips := inst.Get("ips")
		if ips.Exists() && ips.IsArray() {
			ips.ForEach(func(_, ipObj gjson.Result) bool {
				if ipObj.Get("version").Int() == 4 {
					ip = ipObj.Get("ip").String()
					return false
				}
				return true
			})
		}
		table.AddRow(id, ref, itype, region, state, image, ip)
		return true
	})
	table.Render()
	return nil
}

var instancesGetCmd = cli.Command{
	Name:      "get",
	Usage:     "Get instance details",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesGet,
	HideHelpCommand: true,
}

func handleInstancesGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/instances/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesLaunchCmd = cli.Command{
	Name:  "launch",
	Usage: "Launch a new instance",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "region", Usage: "Region (e.g., eu-west-3)", Required: true},
		&cli.StringFlag{Name: "type", Usage: "Instance type (e.g., lsw.c3.large)", Required: true},
		&cli.StringFlag{Name: "image", Usage: "Image ID (e.g., UBUNTU_22_04_64BIT)", Required: true},
		&cli.StringFlag{Name: "contract-type", Usage: "HOURLY or MONTHLY", Required: true},
		&cli.IntFlag{Name: "root-disk-size", Usage: "Root disk size in GB", Required: true},
		&cli.StringFlag{Name: "root-disk-storage-type", Usage: "LOCAL or CENTRAL", Required: true},
		&cli.StringFlag{Name: "reference", Usage: "Instance reference name"},
		&cli.IntFlag{Name: "contract-term", Usage: "Contract term (0,1,3,6,12,24,36)"},
		&cli.IntFlag{Name: "billing-frequency", Usage: "Billing frequency in months (1,3,6,12)"},
		&cli.StringFlag{Name: "ssh-key", Usage: "Public SSH key"},
	},
	Action:          handleInstancesLaunch,
	HideHelpCommand: true,
}

func handleInstancesLaunch(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"region":              cmd.String("region"),
		"type":                cmd.String("type"),
		"imageId":             cmd.String("image"),
		"contractType":        cmd.String("contract-type"),
		"rootDiskSize":        cmd.Int("root-disk-size"),
		"rootDiskStorageType": cmd.String("root-disk-storage-type"),
	}
	if ref := cmd.String("reference"); ref != "" {
		payload["reference"] = ref
	}
	if ct := cmd.Int("contract-term"); ct > 0 {
		payload["contractTerm"] = ct
	}
	if bf := cmd.Int("billing-frequency"); bf > 0 {
		payload["billingFrequency"] = bf
	}
	if sshKey := cmd.String("ssh-key"); sshKey != "" {
		payload["sshKey"] = sshKey
	}

	body, _ := json.Marshal(payload)
	res, err := client.PostJSON(ctx, "/publicCloud/v1/instances", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesTerminateCmd = cli.Command{
	Name:      "terminate",
	Usage:     "Terminate an instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesTerminate,
	HideHelpCommand: true,
}

func handleInstancesTerminate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/publicCloud/v1/instances/"+args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Terminated instance %s\n", args[0])
	return nil
}

var instancesStartCmd = cli.Command{
	Name:      "start",
	Usage:     "Start a stopped instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesStart,
	HideHelpCommand: true,
}

func handleInstancesStart(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/instances/"+args[0]+"/start", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Started instance %s\n", args[0])
	return nil
}

var instancesStopCmd = cli.Command{
	Name:      "stop",
	Usage:     "Stop a running instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesStop,
	HideHelpCommand: true,
}

func handleInstancesStop(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/instances/"+args[0]+"/stop", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Stopped instance %s\n", args[0])
	return nil
}

var instancesRebootCmd = cli.Command{
	Name:      "reboot",
	Usage:     "Reboot an instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesReboot,
	HideHelpCommand: true,
}

func handleInstancesReboot(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/instances/"+args[0]+"/reboot", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Rebooted instance %s\n", args[0])
	return nil
}

var instancesUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "Update an instance",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "New instance type"},
		&cli.StringFlag{Name: "reference", Usage: "New reference"},
		&cli.IntFlag{Name: "root-disk-size", Usage: "New root disk size"},
	},
	Action:          handleInstancesUpdate,
	HideHelpCommand: true,
}

func handleInstancesUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
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
	if rds := cmd.Int("root-disk-size"); rds > 0 {
		payload["rootDiskSize"] = rds
	}

	body, _ := json.Marshal(payload)
	res, err := client.PutJSON(ctx, "/publicCloud/v1/instances/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesConsoleCmd = cli.Command{
	Name:      "console",
	Usage:     "Get console access URL",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesConsole,
	HideHelpCommand: true,
}

func handleInstancesConsole(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/instances/"+args[0]+"/console")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesCredentialsCmd = cli.Command{
	Name:      "credentials",
	Usage:     "List instance credentials",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "Credential type filter"},
	},
	Action:          handleInstancesCredentials,
	HideHelpCommand: true,
}

func handleInstancesCredentials(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/publicCloud/v1/instances/%s/credentials", args[0])
	if t := cmd.String("type"); t != "" {
		path += "/" + t
	}

	res, err := client.Get(ctx, path)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesIPsCmd = cli.Command{
	Name:      "ips",
	Usage:     "List IPs for an instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesIPs,
	HideHelpCommand: true,
}

func handleInstancesIPs(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/instances/"+args[0]+"/ips")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesSnapshotsListCmd = cli.Command{
	Name:      "snapshots",
	Usage:     "List snapshots for an instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesSnapshotsList,
	HideHelpCommand: true,
}

func handleInstancesSnapshotsList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/instances/"+args[0]+"/snapshots")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesSnapshotCreateCmd = cli.Command{
	Name:      "snapshot-create",
	Usage:     "Create snapshot of an instance",
	ArgsUsage: "<instance-id>",
	Action:    handleInstancesSnapshotCreate,
	HideHelpCommand: true,
}

func handleInstancesSnapshotCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, "/publicCloud/v1/instances/"+args[0]+"/snapshots", "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesMetricsCmd = cli.Command{
	Name:      "metrics",
	Usage:     "Get instance metrics (cpu or datatraffic)",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "Metric type: cpu or datatraffic", Value: "cpu"},
		&cli.StringFlag{Name: "from", Usage: "Start date (ISO 8601)", Required: true},
		&cli.StringFlag{Name: "to", Usage: "End date (ISO 8601)", Required: true},
		&cli.StringFlag{Name: "granularity", Usage: "Granularity (5m, 10m, 30m, 1h, etc.)", Value: "1h"},
	},
	Action:          handleInstancesMetrics,
	HideHelpCommand: true,
}

func handleInstancesMetrics(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	metricType := cmd.String("type")
	q := BuildQueryString(map[string]string{
		"from":        cmd.String("from"),
		"to":          cmd.String("to"),
		"granularity": cmd.String("granularity"),
	})
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/metrics/%s%s", args[0], metricType, q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesRegionsCmd = cli.Command{
	Name:   "regions",
	Usage:  "List available regions",
	Action: handleInstancesRegions,
	HideHelpCommand: true,
}

func handleInstancesRegions(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/regions")
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	regions := res.Get("regions")
	if !regions.Exists() || len(regions.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No regions found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "NAME", "LOCATION")
	regions.ForEach(func(_, r gjson.Result) bool {
		table.AddRow(r.Get("name").String(), r.Get("location").String())
		return true
	})
	table.Render()
	return nil
}

var instancesTypesCmd = cli.Command{
	Name:  "types",
	Usage: "List available instance types",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "region", Usage: "Region name", Required: true},
	},
	Action:          handleInstancesTypes,
	HideHelpCommand: true,
}

func handleInstancesTypes(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/instanceTypes?region="+cmd.String("region"))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	types := res.Get("instanceTypes")
	if !types.Exists() || len(types.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No instance types found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "NAME", "CPU", "MEMORY", "HOURLY", "MONTHLY")
	types.ForEach(func(_, t gjson.Result) bool {
		table.AddRow(
			t.Get("name").String(),
			fmt.Sprintf("%d vCPU", t.Get("resources.cpu.value").Int()),
			fmt.Sprintf("%v GiB", t.Get("resources.memory.value").Value()),
			t.Get("prices.hourly").String(),
			t.Get("prices.monthly").String(),
		)
		return true
	})
	table.Render()
	return nil
}

var instancesImagesCmd = cli.Command{
	Name:   "images",
	Usage:  "List available images",
	Action: handleInstancesImages,
	HideHelpCommand: true,
}

func handleInstancesImages(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/images")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
