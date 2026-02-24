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
		&instancesCredentialStoreCmd,
		&instancesCredentialDeleteAllCmd,
		&instancesCredentialGetCmd,
		&instancesCredentialUpdateCmd,
		&instancesCredentialDeleteCmd,
		&instancesIPsCmd,
		&instancesIPGetCmd,
		&instancesIPUpdateCmd,
		&instancesIPNullCmd,
		&instancesIPUnnullCmd,
		&instancesSnapshotsListCmd,
		&instancesSnapshotCreateCmd,
		&instancesSnapshotGetCmd,
		&instancesSnapshotRestoreCmd,
		&instancesSnapshotDeleteCmd,
		&instancesReinstallCmd,
		&instancesReinstallImagesCmd,
		&instancesCancelTerminationCmd,
		&instancesResetPasswordCmd,
		&instancesSecurityGroupsCmd,
		&instancesAttachSecurityGroupsCmd,
		&instancesDetachSecurityGroupsCmd,
		&instancesAddToPrivateNetworkCmd,
		&instancesRemoveFromPrivateNetworkCmd,
		&instancesAttachISOCmd,
		&instancesDetachISOCmd,
		&instancesUserDataCmd,
		&instancesMonitoringEnableCmd,
		&instancesMonitoringStatusCmd,
		&instancesNotifDatatrafficListCmd,
		&instancesNotifDatatrafficGetCmd,
		&instancesNotifDatatrafficCreateCmd,
		&instancesNotifDatatrafficUpdateCmd,
		&instancesNotifDatatrafficDeleteCmd,
		&instancesTypesUpdateCmd,
		&instancesMetricsCmd,
		&instancesRegionsCmd,
		&instancesTypesCmd,
		&instancesImagesCmd,
		&instancesImageCreateCmd,
		&instancesImageUpdateCmd,
		&instancesISOsCmd,
		&instancesExpensesCmd,
		&instancesMarketAppsCmd,
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
	Name:            "images",
	Usage:           "List available images",
	Action:          handleInstancesImages,
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

var instancesCredentialStoreCmd = cli.Command{
	Name:      "credential-store",
	Usage:     "Store credentials for an instance",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "Credential type", Required: true},
		&cli.StringFlag{Name: "username", Usage: "Username", Required: true},
		&cli.StringFlag{Name: "password", Usage: "Password", Required: true},
	},
	Action:          handleInstancesCredentialStore,
	HideHelpCommand: true,
}

func handleInstancesCredentialStore(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
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
	res, err := client.PostJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/credentials", args[0]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesCredentialDeleteAllCmd = cli.Command{
	Name:            "credential-delete-all",
	Usage:           "Delete all credentials for an instance",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesCredentialDeleteAll,
	HideHelpCommand: true,
}

func handleInstancesCredentialDeleteAll(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/credentials", args[0]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted all credentials for %s\n", args[0])
	return nil
}

var instancesCredentialGetCmd = cli.Command{
	Name:            "credential-get",
	Usage:           "Get credentials by type and username",
	ArgsUsage:       "<instance-id> <type> [username]",
	Action:          handleInstancesCredentialGet,
	HideHelpCommand: true,
}

func handleInstancesCredentialGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/publicCloud/v1/instances/%s/credentials/%s", args[0], args[1])
	if len(args) >= 3 {
		path += "/" + args[2]
	}
	res, err := client.Get(ctx, path)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesCredentialUpdateCmd = cli.Command{
	Name:      "credential-update",
	Usage:     "Update credentials for a given type and username",
	ArgsUsage: "<instance-id> <type> <username>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "password", Usage: "New password", Required: true},
	},
	Action:          handleInstancesCredentialUpdate,
	HideHelpCommand: true,
}

func handleInstancesCredentialUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("instance ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"password": cmd.String("password")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/credentials/%s/%s", args[0], args[1], args[2]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesCredentialDeleteCmd = cli.Command{
	Name:            "credential-delete",
	Usage:           "Delete a credential for a given type and username",
	ArgsUsage:       "<instance-id> <type> <username>",
	Action:          handleInstancesCredentialDelete,
	HideHelpCommand: true,
}

func handleInstancesCredentialDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("instance ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted credential %s/%s for %s\n", args[1], args[2], args[0])
	return nil
}

var instancesIPGetCmd = cli.Command{
	Name:            "ip-get",
	Usage:           "Get IP details for an instance",
	ArgsUsage:       "<instance-id> <ip>",
	Action:          handleInstancesIPGet,
	HideHelpCommand: true,
}

func handleInstancesIPGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/ips/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesIPUpdateCmd = cli.Command{
	Name:      "ip-update",
	Usage:     "Update IP address for an instance",
	ArgsUsage: "<instance-id> <ip>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reverse-lookup", Usage: "Reverse lookup hostname", Required: true},
	},
	Action:          handleInstancesIPUpdate,
	HideHelpCommand: true,
}

func handleInstancesIPUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/ips/%s", args[0], args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesIPNullCmd = cli.Command{
	Name:            "ip-null",
	Usage:           "Null route an IP for an instance",
	ArgsUsage:       "<instance-id> <ip>",
	Action:          handleInstancesIPNull,
	HideHelpCommand: true,
}

func handleInstancesIPNull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/ips/%s/null", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s on %s\n", args[1], args[0])
	return nil
}

var instancesIPUnnullCmd = cli.Command{
	Name:            "ip-unnull",
	Usage:           "Remove null route from an IP for an instance",
	ArgsUsage:       "<instance-id> <ip>",
	Action:          handleInstancesIPUnnull,
	HideHelpCommand: true,
}

func handleInstancesIPUnnull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/ips/%s/unnull", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s on %s\n", args[1], args[0])
	return nil
}

var instancesSnapshotGetCmd = cli.Command{
	Name:            "snapshot-get",
	Usage:           "Get snapshot details",
	ArgsUsage:       "<instance-id> <snapshot-id>",
	Action:          handleInstancesSnapshotGet,
	HideHelpCommand: true,
}

func handleInstancesSnapshotGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/snapshots/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesSnapshotRestoreCmd = cli.Command{
	Name:            "snapshot-restore",
	Usage:           "Restore an instance snapshot",
	ArgsUsage:       "<instance-id> <snapshot-id>",
	Action:          handleInstancesSnapshotRestore,
	HideHelpCommand: true,
}

func handleInstancesSnapshotRestore(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/snapshots/%s", args[0], args[1]), []byte("{}"))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Restored snapshot %s for %s\n", args[1], args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesSnapshotDeleteCmd = cli.Command{
	Name:            "snapshot-delete",
	Usage:           "Delete an instance snapshot",
	ArgsUsage:       "<instance-id> <snapshot-id>",
	Action:          handleInstancesSnapshotDelete,
	HideHelpCommand: true,
}

func handleInstancesSnapshotDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/snapshots/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted snapshot %s for %s\n", args[1], args[0])
	return nil
}

var instancesReinstallCmd = cli.Command{
	Name:      "reinstall",
	Usage:     "Reinstall an instance",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "image", Usage: "Image ID for reinstall", Required: true},
	},
	Action:          handleInstancesReinstall,
	HideHelpCommand: true,
}

func handleInstancesReinstall(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"imageId": cmd.String("image")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/reinstall", args[0]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesReinstallImagesCmd = cli.Command{
	Name:            "reinstall-images",
	Usage:           "List images available for reinstall",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesReinstallImages,
	HideHelpCommand: true,
}

func handleInstancesReinstallImages(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/reinstall/images", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesCancelTerminationCmd = cli.Command{
	Name:            "cancel-termination",
	Usage:           "Cancel instance termination",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesCancelTermination,
	HideHelpCommand: true,
}

func handleInstancesCancelTermination(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/cancelTermination", args[0]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Cancelled termination for %s\n", args[0])
	return nil
}

var instancesResetPasswordCmd = cli.Command{
	Name:            "reset-password",
	Usage:           "Reset password for an instance",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesResetPassword,
	HideHelpCommand: true,
}

func handleInstancesResetPassword(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/resetPassword", args[0]), "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesSecurityGroupsCmd = cli.Command{
	Name:            "security-groups",
	Usage:           "Get instance security groups",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesSecurityGroups,
	HideHelpCommand: true,
}

func handleInstancesSecurityGroups(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/securityGroups", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesAttachSecurityGroupsCmd = cli.Command{
	Name:      "attach-security-groups",
	Usage:     "Attach security groups to instance",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload with security group IDs", Required: true},
	},
	Action:          handleInstancesAttachSecurityGroups,
	HideHelpCommand: true,
}

func handleInstancesAttachSecurityGroups(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/attachSecurityGroups", args[0]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Attached security groups to %s\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesDetachSecurityGroupsCmd = cli.Command{
	Name:      "detach-security-groups",
	Usage:     "Detach security groups from instance",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload with security group IDs", Required: true},
	},
	Action:          handleInstancesDetachSecurityGroups,
	HideHelpCommand: true,
}

func handleInstancesDetachSecurityGroups(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/detachSecurityGroups", args[0]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Detached security groups from %s\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesAddToPrivateNetworkCmd = cli.Command{
	Name:      "add-to-private-network",
	Usage:     "Add instance to private network",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "private-network-id", Usage: "Private network ID", Required: true},
	},
	Action:          handleInstancesAddToPrivateNetwork,
	HideHelpCommand: true,
}

func handleInstancesAddToPrivateNetwork(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"privateNetworkId": cmd.String("private-network-id")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/addToPrivateNetwork", args[0]), body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Added %s to private network\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesRemoveFromPrivateNetworkCmd = cli.Command{
	Name:            "remove-from-private-network",
	Usage:           "Remove instance from private network",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesRemoveFromPrivateNetwork,
	HideHelpCommand: true,
}

func handleInstancesRemoveFromPrivateNetwork(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/removeFromPrivateNetwork", args[0]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed %s from private network\n", args[0])
	return nil
}

var instancesAttachISOCmd = cli.Command{
	Name:      "attach-iso",
	Usage:     "Attach ISO to an instance",
	ArgsUsage: "<instance-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "iso-id", Usage: "ISO ID", Required: true},
	},
	Action:          handleInstancesAttachISO,
	HideHelpCommand: true,
}

func handleInstancesAttachISO(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"isoId": cmd.String("iso-id")})
	res, err := client.PostJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/attachIso", args[0]), body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Attached ISO to %s\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesDetachISOCmd = cli.Command{
	Name:            "detach-iso",
	Usage:           "Detach ISO from an instance",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesDetachISO,
	HideHelpCommand: true,
}

func handleInstancesDetachISO(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/detachIso", args[0]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Detached ISO from %s\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesUserDataCmd = cli.Command{
	Name:            "user-data",
	Usage:           "Get user data for an instance",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesUserData,
	HideHelpCommand: true,
}

func handleInstancesUserData(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/userData", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesMonitoringEnableCmd = cli.Command{
	Name:            "monitoring-enable",
	Usage:           "Enable monitoring for an instance",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesMonitoringEnable,
	HideHelpCommand: true,
}

func handleInstancesMonitoringEnable(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/instance/%s/monitoring/enable", args[0]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Enabled monitoring for %s\n", args[0])
	return nil
}

var instancesMonitoringStatusCmd = cli.Command{
	Name:            "monitoring-status",
	Usage:           "Get monitoring status for an instance",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesMonitoringStatus,
	HideHelpCommand: true,
}

func handleInstancesMonitoringStatus(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instance/%s/monitoring/status", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesNotifDatatrafficListCmd = cli.Command{
	Name:            "notif-datatraffic-list",
	Usage:           "List data traffic notification settings",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesNotifDatatrafficList,
	HideHelpCommand: true,
}

func handleInstancesNotifDatatrafficList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/notificationSettings/dataTraffic", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesNotifDatatrafficGetCmd = cli.Command{
	Name:            "notif-datatraffic-get",
	Usage:           "Get a data traffic notification setting",
	ArgsUsage:       "<instance-id> <notification-id>",
	Action:          handleInstancesNotifDatatrafficGet,
	HideHelpCommand: true,
}

func handleInstancesNotifDatatrafficGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/notificationSettings/dataTraffic/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesNotifDatatrafficCreateCmd = cli.Command{
	Name:      "notif-datatraffic-create",
	Usage:     "Create a data traffic notification setting",
	ArgsUsage: "<instance-id> <notification-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleInstancesNotifDatatrafficCreate,
	HideHelpCommand: true,
}

func handleInstancesNotifDatatrafficCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/notificationSettings/dataTraffic/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesNotifDatatrafficUpdateCmd = cli.Command{
	Name:      "notif-datatraffic-update",
	Usage:     "Update a data traffic notification setting",
	ArgsUsage: "<instance-id> <notification-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleInstancesNotifDatatrafficUpdate,
	HideHelpCommand: true,
}

func handleInstancesNotifDatatrafficUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/notificationSettings/dataTraffic/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesNotifDatatrafficDeleteCmd = cli.Command{
	Name:            "notif-datatraffic-delete",
	Usage:           "Delete a data traffic notification setting",
	ArgsUsage:       "<instance-id> <notification-id>",
	Action:          handleInstancesNotifDatatrafficDelete,
	HideHelpCommand: true,
}

func handleInstancesNotifDatatrafficDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("instance ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/notificationSettings/dataTraffic/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted notification %s for %s\n", args[1], args[0])
	return nil
}

var instancesTypesUpdateCmd = cli.Command{
	Name:            "types-update",
	Usage:           "List available instance types for update",
	ArgsUsage:       "<instance-id>",
	Action:          handleInstancesTypesUpdate,
	HideHelpCommand: true,
}

func handleInstancesTypesUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("instance ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/instances/%s/instanceTypesUpdate", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesImageCreateCmd = cli.Command{
	Name:  "image-create",
	Usage: "Create a custom image",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload for image creation", Required: true},
	},
	Action:          handleInstancesImageCreate,
	HideHelpCommand: true,
}

func handleInstancesImageCreate(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/publicCloud/v1/images", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesImageUpdateCmd = cli.Command{
	Name:      "image-update",
	Usage:     "Update a custom image",
	ArgsUsage: "<image-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload for image update", Required: true},
	},
	Action:          handleInstancesImageUpdate,
	HideHelpCommand: true,
}

func handleInstancesImageUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("image ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/publicCloud/v1/images/"+args[0], []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesISOsCmd = cli.Command{
	Name:            "isos",
	Usage:           "List available ISOs",
	Action:          handleInstancesISOs,
	HideHelpCommand: true,
}

func handleInstancesISOs(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/isos")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesExpensesCmd = cli.Command{
	Name:      "expenses",
	Usage:     "Get costs for an equipment",
	ArgsUsage: "<equipment-id>",
	Action:    handleInstancesExpenses,
	HideHelpCommand: true,
}

func handleInstancesExpenses(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("equipment ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/equipments/%s/expenses", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var instancesMarketAppsCmd = cli.Command{
	Name:            "market-apps",
	Usage:           "List marketplace apps",
	Action:          handleInstancesMarketApps,
	HideHelpCommand: true,
}

func handleInstancesMarketApps(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/marketApps")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}
