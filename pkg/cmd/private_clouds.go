package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var privateCloudsCmd = cli.Command{
	Name:    "private-clouds",
	Aliases: []string{"pc"},
	Usage:   "Manage private clouds",
	Commands: []*cli.Command{
		&pcListCmd, &pcGetCmd,
		&pcCredentialsByTypeCmd, &pcCredentialGetCmd,
		&pcMetricsBandwidthCmd, &pcMetricsCPUCmd, &pcMetricsMemoryCmd,
		&pcMetricsStorageCmd, &pcMetricsDatatrafficCmd,
	},
	HideHelpCommand: true,
}

var pcListCmd = cli.Command{Name: "list", Usage: "List private clouds", Flags: PaginationFlags, Action: handlePCList, HideHelpCommand: true}

func handlePCList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cloud/v2/privateClouds?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var pcGetCmd = cli.Command{Name: "get", Usage: "Inspect a private cloud", ArgsUsage: "<id>", Action: handlePCGet, HideHelpCommand: true}

func handlePCGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("private cloud ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cloud/v2/privateClouds/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var pcCredentialsByTypeCmd = cli.Command{Name: "credentials", Usage: "List credentials by type", ArgsUsage: "<id> <type>", Action: handlePCCredentials, HideHelpCommand: true}

func handlePCCredentials(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("private cloud ID and credential type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/privateClouds/%s/credentials/%s", args[0], args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var pcCredentialGetCmd = cli.Command{Name: "credential-get", Usage: "Get credential", ArgsUsage: "<id> <type> <username>", Action: handlePCCredentialGet, HideHelpCommand: true}

func handlePCCredentialGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("private cloud ID, type, and username required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/privateClouds/%s/credentials/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

func pcMetricsHandler(metricType string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < 1 {
			return fmt.Errorf("private cloud ID required")
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "granularity": cmd.String("granularity")})
		res, err := client.Get(ctx, fmt.Sprintf("/cloud/v2/privateClouds/%s/metrics/%s%s", args[0], metricType, q))
		if err != nil {
			return err
		}
		return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
	}
}

var pcMetricsFlags = []cli.Flag{
	&cli.StringFlag{Name: "from", Required: true},
	&cli.StringFlag{Name: "to", Required: true},
	&cli.StringFlag{Name: "granularity", Value: "1h"},
}

var pcMetricsBandwidthCmd = cli.Command{Name: "metrics-bandwidth", Usage: "Bandwidth metrics", ArgsUsage: "<id>", Flags: pcMetricsFlags, Action: pcMetricsHandler("bandwidth"), HideHelpCommand: true}
var pcMetricsCPUCmd = cli.Command{Name: "metrics-cpu", Usage: "CPU metrics", ArgsUsage: "<id>", Flags: pcMetricsFlags, Action: pcMetricsHandler("cpu"), HideHelpCommand: true}
var pcMetricsMemoryCmd = cli.Command{Name: "metrics-memory", Usage: "Memory metrics", ArgsUsage: "<id>", Flags: pcMetricsFlags, Action: pcMetricsHandler("memory"), HideHelpCommand: true}
var pcMetricsStorageCmd = cli.Command{Name: "metrics-storage", Usage: "Storage metrics", ArgsUsage: "<id>", Flags: pcMetricsFlags, Action: pcMetricsHandler("storage"), HideHelpCommand: true}
var pcMetricsDatatrafficCmd = cli.Command{Name: "metrics-datatraffic", Usage: "Datatraffic metrics", ArgsUsage: "<id>", Flags: pcMetricsFlags, Action: pcMetricsHandler("datatraffic"), HideHelpCommand: true}
