package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var virtualServersCmd = cli.Command{
	Name:    "virtual-servers",
	Aliases: []string{"vs"},
	Usage:   "Manage virtual servers",
	Commands: []*cli.Command{
		&vsListCmd, &vsGetCmd, &vsUpdateCmd,
		&vsPowerOnCmd, &vsPowerOffCmd, &vsRebootCmd, &vsReinstallCmd,
		&vsCredentialsUpdateCmd, &vsCredentialsByTypeCmd, &vsCredentialGetCmd,
		&vsMetricsCmd,
		&vsSnapshotsCmd, &vsSnapshotCreateCmd, &vsSnapshotGetCmd, &vsSnapshotDeleteCmd, &vsSnapshotRestoreCmd,
		&vsTemplatesCmd,
	},
	HideHelpCommand: true,
}

var vsListCmd = cli.Command{Name: "list", Usage: "List virtual servers", Flags: PaginationFlags, Action: handleVSList, HideHelpCommand: true}

func handleVSList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cloud/v2/virtualServers?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsGetCmd = cli.Command{Name: "get", Usage: "Get virtual server details", ArgsUsage: "<id>", Action: handleVSGet, HideHelpCommand: true}

func handleVSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cloud/v2/virtualServers/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsUpdateCmd = cli.Command{
	Name: "update", Usage: "Update a virtual server", ArgsUsage: "<id>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true}},
	Action: handleVSUpdate, HideHelpCommand: true,
}

func handleVSUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/cloud/v2/virtualServers/"+args[0], []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsPowerOnCmd = cli.Command{Name: "power-on", Usage: "Power on", ArgsUsage: "<id>", Action: handleVSPowerOn, HideHelpCommand: true}

func handleVSPowerOn(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/cloud/v2/virtualServers/"+args[0]+"/powerOn", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Powered on %s\n", args[0])
	return nil
}

var vsPowerOffCmd = cli.Command{Name: "power-off", Usage: "Power off", ArgsUsage: "<id>", Action: handleVSPowerOff, HideHelpCommand: true}

func handleVSPowerOff(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/cloud/v2/virtualServers/"+args[0]+"/powerOff", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Powered off %s\n", args[0])
	return nil
}

var vsRebootCmd = cli.Command{Name: "reboot", Usage: "Reboot", ArgsUsage: "<id>", Action: handleVSReboot, HideHelpCommand: true}

func handleVSReboot(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Post(ctx, "/cloud/v2/virtualServers/"+args[0]+"/reboot", "")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Rebooted %s\n", args[0])
	return nil
}

var vsReinstallCmd = cli.Command{
	Name: "reinstall", Usage: "Reinstall", ArgsUsage: "<id>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true}},
	Action: handleVSReinstall, HideHelpCommand: true,
}

func handleVSReinstall(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/cloud/v2/virtualServers/"+args[0]+"/reinstall", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsCredentialsUpdateCmd = cli.Command{
	Name: "credentials-update", Usage: "Update credentials", ArgsUsage: "<id>",
	Flags:  []cli.Flag{&cli.StringFlag{Name: "payload", Required: true}},
	Action: handleVSCredentialsUpdate, HideHelpCommand: true,
}

func handleVSCredentialsUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/cloud/v2/virtualServers/"+args[0]+"/credentials", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsCredentialsByTypeCmd = cli.Command{Name: "credentials", Usage: "List credentials by type", ArgsUsage: "<id> <type>", Action: handleVSCredentialsByType, HideHelpCommand: true}

func handleVSCredentialsByType(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("virtual server ID and credential type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/virtualServers/%s/credentials/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsCredentialGetCmd = cli.Command{Name: "credential-get", Usage: "Get credential", ArgsUsage: "<id> <type> <username>", Action: handleVSCredentialGet, HideHelpCommand: true}

func handleVSCredentialGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("virtual server ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/virtualServers/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsMetricsCmd = cli.Command{
	Name: "metrics", Usage: "Get datatraffic metrics", ArgsUsage: "<id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true},
		&cli.StringFlag{Name: "granularity", Value: "1h"},
	},
	Action: handleVSMetrics, HideHelpCommand: true,
}

func handleVSMetrics(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "granularity": cmd.String("granularity")})
	res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/virtualServers/%s/metrics/datatraffic%s", args[0], q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsSnapshotsCmd = cli.Command{Name: "snapshots", Usage: "List snapshots", ArgsUsage: "<id>", Action: handleVSSnapshots, HideHelpCommand: true}

func handleVSSnapshots(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cloud/v2/virtualServers/"+args[0]+"/snapshots")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsSnapshotCreateCmd = cli.Command{Name: "snapshot-create", Usage: "Create snapshot", ArgsUsage: "<id>", Action: handleVSSnapshotCreate, HideHelpCommand: true}

func handleVSSnapshotCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, "/cloud/v2/virtualServers/"+args[0]+"/snapshots", "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsSnapshotGetCmd = cli.Command{Name: "snapshot-get", Usage: "Get snapshot", ArgsUsage: "<id> <snapshot-id>", Action: handleVSSnapshotGet, HideHelpCommand: true}

func handleVSSnapshotGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("virtual server ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/virtualServers/%s/snapshots/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsSnapshotDeleteCmd = cli.Command{Name: "snapshot-delete", Usage: "Delete snapshot", ArgsUsage: "<id> <snapshot-id>", Action: handleVSSnapshotDelete, HideHelpCommand: true}

func handleVSSnapshotDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("virtual server ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/cloud/v2/virtualServers/%s/snapshots/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted snapshot %s\n", args[1])
	return nil
}

var vsSnapshotRestoreCmd = cli.Command{Name: "snapshot-restore", Usage: "Restore snapshot", ArgsUsage: "<id> <snapshot-id>", Action: handleVSSnapshotRestore, HideHelpCommand: true}

func handleVSSnapshotRestore(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("virtual server ID and snapshot ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, fmt.Sprintf("/cloud/v2/virtualServers/%s/snapshots/%s/restore", args[0], args[1]), "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var vsTemplatesCmd = cli.Command{Name: "templates", Usage: "List templates", ArgsUsage: "<id>", Action: handleVSTemplates, HideHelpCommand: true}

func handleVSTemplates(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("virtual server ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cloud/v2/virtualServers/"+args[0]+"/templates")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
