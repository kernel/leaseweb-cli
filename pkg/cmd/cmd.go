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
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "Leaseweb API key (overrides profile)",
			},
			&cli.StringFlag{
				Name:  "base-url",
				Usage: "Override the base URL for API requests",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging of HTTP requests",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
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
			&configCmd,
			&dedicatedServersCmd,
			&domainsCmd,
			&instancesCmd,
			&invoicesCmd,
			&ipsCmd,
			&loadBalancersCmd,
			&privateNetworksCmd,
			&servicesCmd,
		},
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "@completion",
		HideHelpCommand:            true,
	}
}
