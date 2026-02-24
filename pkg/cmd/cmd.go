package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

var Command *cli.Command

func init() {
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Fprintf(os.Stdout, "lw version %s\n", cmd.Root().Version)
	}
	Command = &cli.Command{
		Name:    "lw",
		Usage:   "CLI for the Leaseweb API",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Usage:   "Config profile to use",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging of HTTP requests",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output format (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value:   "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("output must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "transform",
				Usage: "GJSON expression to transform output",
			},
		},
		Commands: []*cli.Command{
			&abuseReportsCmd,
			&acronisBackupCmd,
			&aggregationPacksCmd,
			&apiKeysCmd,
			&cdnCmd,
			&colocationsCmd,
			&configCmd,
			&datacenterAccessCmd,
			&dedicatedRacksCmd,
			&dedicatedServersCmd,
			&domainsCmd,
			&emailsCmd,
			&floatingIPsCmd,
			&instancesCmd,
			&invoicesCmd,
			&ipsCmd,
			&loadBalancersCmd,
			&networkEquipmentCmd,
			&privateCloudsCmd,
			&privateNetworksCmd,
			&remoteManagementCmd,
			&servicesCmd,
			&storageCmd,
			&trafficPolicyCmd,
			&virtualServersCmd,
			&vpsCmd,
			&webhostingCmd,
		},
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "@completion",
		HideHelpCommand:            true,
	}
}
