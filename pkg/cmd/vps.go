package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var vpsCmd = cli.Command{
	Name:    "vps",
	Aliases: []string{"v"},
	Usage:   "Manage VPS instances",
	Commands: []*cli.Command{
		&vpsListCmd, &vpsGetCmd, &vpsUpdateCmd,
		&vpsStartCmd, &vpsStopCmd, &vpsRebootCmd,
		&vpsReinstallCmd, &vpsReinstallImagesCmd, &vpsResetPasswordCmd,
		&vpsConsoleCmd,
		&vpsCredentialsCmd, &vpsCredentialStoreCmd, &vpsCredentialDeleteAllCmd,
		&vpsCredentialGetCmd, &vpsCredentialUpdateCmd, &vpsCredentialDeleteCmd,
		&vpsIPsCmd, &vpsIPGetCmd, &vpsIPUpdateCmd, &vpsIPNullCmd, &vpsIPUnnullCmd,
		&vpsSnapshotsCmd, &vpsSnapshotCreateCmd, &vpsSnapshotGetCmd,
		&vpsSnapshotRestoreCmd, &vpsSnapshotDeleteCmd,
		&vpsMetricsCmd,
		&vpsMonitoringEnableCmd, &vpsMonitoringStatusCmd,
		&vpsNotifDatatrafficListCmd, &vpsNotifDatatrafficGetCmd,
		&vpsNotifDatatrafficCreateCmd, &vpsNotifDatatrafficUpdateCmd, &vpsNotifDatatrafficDeleteCmd,
		&vpsAttachISOCmd, &vpsDetachISOCmd, &vpsISOsCmd,
	},
	HideHelpCommand: true,
}

var vpsListCmd = cli.Command{
	Name: "list", Usage: "List VPS instances", Flags: PaginationFlags,
	Action: handleVPSList, HideHelpCommand: true,
}

func handleVPSList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	format := cmd.Root().String("output")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}
	vps := res.Get("vps")
	if !vps.Exists() || len(vps.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No VPS found.")
		return nil
	}
	table := NewTableWriter(os.Stdout, "ID", "REFERENCE", "STATE", "REGION", "IP")
	table.TruncOrder = []int{0}
	vps.ForEach(func(_, v gjson.Result) bool {
		ip := ""
		v.Get("ips").ForEach(func(_, ipObj gjson.Result) bool {
			if ipObj.Get("version").Int() == 4 {
				ip = ipObj.Get("ip").String()
				return false
			}
			return true
		})
		table.AddRow(v.Get("id").String(), v.Get("reference").String(), v.Get("state").String(), v.Get("region").String(), ip)
		return true
	})
	table.Render()
	return nil
}

var vpsGetCmd = cli.Command{
	Name: "get", Usage: "Get VPS details", ArgsUsage: "<vps-id>",
	Action: handleVPSGet, HideHelpCommand: true,
}

func handleVPSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsUpdateCmd = cli.Command{
	Name: "update", Usage: "Update VPS details", ArgsUsage: "<vps-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "reference", Usage: "New reference"},
	},
	Action: handleVPSUpdate, HideHelpCommand: true,
}

func handleVPSUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	payload := map[string]any{}
	if ref := cmd.String("reference"); ref != "" {
		payload["reference"] = ref
	}
	body, _ := json.Marshal(payload)
	res, err := client.PutJSON(ctx, "/publicCloud/v1/vps/"+args[0], body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsStartCmd = cli.Command{Name: "start", Usage: "Start a VPS", ArgsUsage: "<vps-id>", Action: handleVPSStart, HideHelpCommand: true}

func handleVPSStart(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/start", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Started VPS %s\n", args[0])
	return nil
}

var vpsStopCmd = cli.Command{Name: "stop", Usage: "Stop a VPS", ArgsUsage: "<vps-id>", Action: handleVPSStop, HideHelpCommand: true}

func handleVPSStop(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/stop", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Stopped VPS %s\n", args[0])
	return nil
}

var vpsRebootCmd = cli.Command{Name: "reboot", Usage: "Reboot a VPS", ArgsUsage: "<vps-id>", Action: handleVPSReboot, HideHelpCommand: true}

func handleVPSReboot(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/reboot", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Rebooted VPS %s\n", args[0])
	return nil
}

var vpsReinstallCmd = cli.Command{
	Name: "reinstall", Usage: "Reinstall a VPS", ArgsUsage: "<vps-id>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "image", Usage: "Image ID", Required: true}},
	Action: handleVPSReinstall, HideHelpCommand: true,
}

func handleVPSReinstall(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"imageId": cmd.String("image")})
	res, err := client.PutJSON(ctx, "/publicCloud/v1/vps/"+args[0]+"/reinstall", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsReinstallImagesCmd = cli.Command{
	Name: "reinstall-images", Usage: "List images available for reinstall", ArgsUsage: "<vps-id>",
	Action: handleVPSReinstallImages, HideHelpCommand: true,
}

func handleVPSReinstallImages(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/reinstall/images")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsResetPasswordCmd = cli.Command{Name: "reset-password", Usage: "Reset password", ArgsUsage: "<vps-id>", Action: handleVPSResetPassword, HideHelpCommand: true}

func handleVPSResetPassword(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/resetPassword", "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsConsoleCmd = cli.Command{Name: "console", Usage: "Get console access", ArgsUsage: "<vps-id>", Action: handleVPSConsole, HideHelpCommand: true}

func handleVPSConsole(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/console")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsCredentialsCmd = cli.Command{Name: "credentials", Usage: "List credentials", ArgsUsage: "<vps-id>", Action: handleVPSCredentials, HideHelpCommand: true}

func handleVPSCredentials(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/credentials")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsCredentialStoreCmd = cli.Command{
	Name: "credential-store", Usage: "Store credentials", ArgsUsage: "<vps-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Required: true}, &cli.StringFlag{Name: "username", Required: true}, &cli.StringFlag{Name: "password", Required: true},
	},
	Action: handleVPSCredentialStore, HideHelpCommand: true,
}

func handleVPSCredentialStore(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"type": cmd.String("type"), "username": cmd.String("username"), "password": cmd.String("password")})
	res, err := client.PostJSON(ctx, "/publicCloud/v1/vps/"+args[0]+"/credentials", body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsCredentialDeleteAllCmd = cli.Command{Name: "credential-delete-all", Usage: "Delete all credentials", ArgsUsage: "<vps-id>", Action: handleVPSCredentialDeleteAll, HideHelpCommand: true}

func handleVPSCredentialDeleteAll(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/publicCloud/v1/vps/"+args[0]+"/credentials")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted all credentials for %s\n", args[0])
	return nil
}

var vpsCredentialGetCmd = cli.Command{Name: "credential-get", Usage: "Get credentials by type and username", ArgsUsage: "<vps-id> <type> [username]", Action: handleVPSCredentialGet, HideHelpCommand: true}

func handleVPSCredentialGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/publicCloud/v1/vps/%s/credentials/%s", args[0], args[1])
	if len(args) >= 3 {
		path += "/" + args[2]
	}
	res, err := client.Get(ctx, path)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsCredentialUpdateCmd = cli.Command{
	Name: "credential-update", Usage: "Update credentials", ArgsUsage: "<vps-id> <type> <username>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "password", Required: true}},
	Action: handleVPSCredentialUpdate, HideHelpCommand: true,
}

func handleVPSCredentialUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("VPS ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"password": cmd.String("password")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/credentials/%s/%s", args[0], args[1], args[2]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsCredentialDeleteCmd = cli.Command{Name: "credential-delete", Usage: "Delete credential", ArgsUsage: "<vps-id> <type> <username>", Action: handleVPSCredentialDelete, HideHelpCommand: true}

func handleVPSCredentialDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("VPS ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted credential %s/%s for %s\n", args[1], args[2], args[0])
	return nil
}

var vpsIPsCmd = cli.Command{Name: "ips", Usage: "List IPs", ArgsUsage: "<vps-id>", Action: handleVPSIPs, HideHelpCommand: true}

func handleVPSIPs(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/ips")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsIPGetCmd = cli.Command{Name: "ip-get", Usage: "Get IP details", ArgsUsage: "<vps-id> <ip>", Action: handleVPSIPGet, HideHelpCommand: true}

func handleVPSIPGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/ips/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsIPUpdateCmd = cli.Command{
	Name: "ip-update", Usage: "Update IP", ArgsUsage: "<vps-id> <ip>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "reverse-lookup", Required: true}},
	Action: handleVPSIPUpdate, HideHelpCommand: true,
}

func handleVPSIPUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"reverseLookup": cmd.String("reverse-lookup")})
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/ips/%s", args[0], args[1]), body)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsIPNullCmd = cli.Command{Name: "ip-null", Usage: "Null route IP", ArgsUsage: "<vps-id> <ip>", Action: handleVPSIPNull, HideHelpCommand: true}

func handleVPSIPNull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/ips/%s/null", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Null routed %s\n", args[1])
	return nil
}

var vpsIPUnnullCmd = cli.Command{Name: "ip-unnull", Usage: "Remove null route", ArgsUsage: "<vps-id> <ip>", Action: handleVPSIPUnnull, HideHelpCommand: true}

func handleVPSIPUnnull(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and IP required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/ips/%s/unnull", args[0], args[1]), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Removed null route from %s\n", args[1])
	return nil
}

var vpsSnapshotsCmd = cli.Command{Name: "snapshots", Usage: "List snapshots", ArgsUsage: "<vps-id>", Action: handleVPSSnapshots, HideHelpCommand: true}

func handleVPSSnapshots(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/snapshots")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsSnapshotCreateCmd = cli.Command{Name: "snapshot-create", Usage: "Create snapshot", ArgsUsage: "<vps-id>", Action: handleVPSSnapshotCreate, HideHelpCommand: true}

func handleVPSSnapshotCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/snapshots", "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsSnapshotGetCmd = cli.Command{Name: "snapshot-get", Usage: "Get snapshot", ArgsUsage: "<vps-id> <snapshot-id>", Action: handleVPSSnapshotGet, HideHelpCommand: true}

func handleVPSSnapshotGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/snapshots/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsSnapshotRestoreCmd = cli.Command{Name: "snapshot-restore", Usage: "Restore snapshot", ArgsUsage: "<vps-id> <snapshot-id>", Action: handleVPSSnapshotRestore, HideHelpCommand: true}

func handleVPSSnapshotRestore(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/snapshots/%s", args[0], args[1]), []byte("{}"))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Restored snapshot %s\n", args[1])
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsSnapshotDeleteCmd = cli.Command{Name: "snapshot-delete", Usage: "Delete snapshot", ArgsUsage: "<vps-id> <snapshot-id>", Action: handleVPSSnapshotDelete, HideHelpCommand: true}

func handleVPSSnapshotDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/snapshots/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted snapshot %s\n", args[1])
	return nil
}

var vpsMetricsCmd = cli.Command{
	Name: "metrics", Usage: "Get VPS metrics", ArgsUsage: "<vps-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true},
		&cli.StringFlag{Name: "granularity", Value: "1h"},
	},
	Action: handleVPSMetrics, HideHelpCommand: true,
}

func handleVPSMetrics(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "granularity": cmd.String("granularity")})
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/metrics/datatraffic%s", args[0], q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsMonitoringEnableCmd = cli.Command{Name: "monitoring-enable", Usage: "Enable monitoring", ArgsUsage: "<vps-id>", Action: handleVPSMonitoringEnable, HideHelpCommand: true}

func handleVPSMonitoringEnable(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/monitoring/enable", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Enabled monitoring for %s\n", args[0])
	return nil
}

var vpsMonitoringStatusCmd = cli.Command{Name: "monitoring-status", Usage: "Get monitoring status", ArgsUsage: "<vps-id>", Action: handleVPSMonitoringStatus, HideHelpCommand: true}

func handleVPSMonitoringStatus(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/monitoring/status")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsNotifDatatrafficListCmd = cli.Command{Name: "notif-datatraffic-list", Usage: "List data traffic notifications", ArgsUsage: "<vps-id>", Action: handleVPSNotifList, HideHelpCommand: true}

func handleVPSNotifList(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/"+args[0]+"/notificationSettings/dataTraffic")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsNotifDatatrafficGetCmd = cli.Command{Name: "notif-datatraffic-get", Usage: "Get notification", ArgsUsage: "<vps-id> <id>", Action: handleVPSNotifGet, HideHelpCommand: true}

func handleVPSNotifGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/notificationSettings/dataTraffic/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsNotifDatatrafficCreateCmd = cli.Command{
	Name: "notif-datatraffic-create", Usage: "Create notification", ArgsUsage: "<vps-id> <id>",
	Flags: []cli.Flag{&cli.StringFlag{Name: "payload", Required: true}}, Action: handleVPSNotifCreate, HideHelpCommand: true,
}

func handleVPSNotifCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/notificationSettings/dataTraffic/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsNotifDatatrafficUpdateCmd = cli.Command{
	Name: "notif-datatraffic-update", Usage: "Update notification", ArgsUsage: "<vps-id> <id>",
	Flags: []cli.Flag{&cli.StringFlag{Name: "payload", Required: true}}, Action: handleVPSNotifUpdate, HideHelpCommand: true,
}

func handleVPSNotifUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/notificationSettings/dataTraffic/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsNotifDatatrafficDeleteCmd = cli.Command{Name: "notif-datatraffic-delete", Usage: "Delete notification", ArgsUsage: "<vps-id> <id>", Action: handleVPSNotifDelete, HideHelpCommand: true}

func handleVPSNotifDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("VPS ID and notification ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/publicCloud/v1/vps/%s/notificationSettings/dataTraffic/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted notification %s\n", args[1])
	return nil
}

var vpsAttachISOCmd = cli.Command{
	Name: "attach-iso", Usage: "Attach ISO", ArgsUsage: "<vps-id>",
	Flags: []cli.Flag{&cli.StringFlag{Name: "iso-id", Required: true}}, Action: handleVPSAttachISO, HideHelpCommand: true,
}

func handleVPSAttachISO(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]string{"isoId": cmd.String("iso-id")})
	res, err := client.PostJSON(ctx, "/publicCloud/v1/vps/"+args[0]+"/attachIso", body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Attached ISO to %s\n", args[0])
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsDetachISOCmd = cli.Command{Name: "detach-iso", Usage: "Detach ISO", ArgsUsage: "<vps-id>", Action: handleVPSDetachISO, HideHelpCommand: true}

func handleVPSDetachISO(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("VPS ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, "/publicCloud/v1/vps/"+args[0]+"/detachIso", "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vpsISOsCmd = cli.Command{Name: "isos", Usage: "List available ISOs", Action: handleVPSISOs, HideHelpCommand: true}

func handleVPSISOs(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/publicCloud/v1/vps/isos")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
