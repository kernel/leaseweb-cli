package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var storageCmd = cli.Command{
	Name:  "storage",
	Usage: "Manage storage",
	Commands: []*cli.Command{
		&storageListCmd,
		&storageListVMsCmd,
		&storageVMJobCmd,
		&storageVolumesCmd,
		&storageVolumeGrowCmd,
	},
	HideHelpCommand: true,
}

var storageListCmd = cli.Command{
	Name:            "list",
	Usage:           "List storages",
	Action:          handleStorageList,
	HideHelpCommand: true,
}

func handleStorageList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/storage/v1/storages")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var storageListVMsCmd = cli.Command{
	Name:            "vms",
	Usage:           "List storage VMs",
	Action:          handleStorageListVMs,
	HideHelpCommand: true,
}

func handleStorageListVMs(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/storage/v1/storageVMs")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var storageVMJobCmd = cli.Command{
	Name:            "vm-job",
	Usage:           "Get storage VM job",
	ArgsUsage:       "<storage-vm-id> <job-id>",
	Action:          handleStorageVMJob,
	HideHelpCommand: true,
}

func handleStorageVMJob(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("storage VM ID and job ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/storage/v1/storageVMs/%s/jobs/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var storageVolumesCmd = cli.Command{
	Name:            "volumes",
	Usage:           "List volumes for a storage VM",
	ArgsUsage:       "<storage-vm-id>",
	Action:          handleStorageVolumes,
	HideHelpCommand: true,
}

func handleStorageVolumes(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("storage VM ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/storage/v1/storageVMs/%s/volumes", args[0]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var storageVolumeGrowCmd = cli.Command{
	Name:      "volume-grow",
	Usage:     "Grow a storage volume",
	ArgsUsage: "<storage-vm-id> <volume-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleStorageVolumeGrow,
	HideHelpCommand: true,
}

func handleStorageVolumeGrow(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("storage VM ID and volume ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, fmt.Sprintf("/storage/v1/storageVMs/%s/volumes/%s/grow", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
