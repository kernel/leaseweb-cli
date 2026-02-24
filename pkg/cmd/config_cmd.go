package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

var configCmd = cli.Command{
	Name:  "config",
	Usage: "Manage CLI configuration",
	Commands: []*cli.Command{
		&configInitCmd,
		&configShowCmd,
	},
	HideHelpCommand: true,
}

var configInitCmd = cli.Command{
	Name:            "init",
	Usage:           "Initialize or update CLI configuration",
	Action:          handleConfigInit,
	HideHelpCommand: true,
}

var configShowCmd = cli.Command{
	Name:            "show",
	Usage:           "Show current configuration",
	Action:          handleConfigShow,
	HideHelpCommand: true,
}

func handleConfigInit(_ context.Context, cmd *cli.Command) error {
	reader := bufio.NewReader(os.Stdin)
	cfg := loadConfig()

	fmt.Println("Leaseweb CLI Configuration")
	fmt.Println()

	if len(cfg.Profiles) > 0 {
		fmt.Printf("Existing profiles: %s\n", strings.Join(profileNames(cfg), ", "))
		fmt.Println()
	}

	for {
		name := prompt(reader, "Profile name (or 'done' to finish)", "default")
		if strings.ToLower(name) == "done" {
			break
		}

		existing := ""
		if p, ok := cfg.Profiles[name]; ok {
			existing = p.APIKey
		}

		apiKey := prompt(reader, fmt.Sprintf("API key for %q", name), existing)
		if apiKey == "" {
			fmt.Println("Skipping profile (no API key provided)")
			continue
		}

		cfg.Profiles[name] = ProfileConfig{APIKey: apiKey}
		fmt.Printf("Profile %q saved.\n", name)

		if cfg.DefaultProfile == "" {
			cfg.DefaultProfile = name
		}
	}

	if len(cfg.Profiles) > 1 || (len(cfg.Profiles) == 1 && cfg.DefaultProfile == "") {
		def := prompt(reader, "Default profile", cfg.DefaultProfile)
		if def != "" {
			cfg.DefaultProfile = def
		}
	}

	if err := writeConfig(cfg); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Printf("\nConfiguration written to %s\n", getConfigPath())
	return nil
}

func handleConfigShow(_ context.Context, cmd *cli.Command) error {
	cfg := loadConfig()

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles configured. Run 'lw config init' to set up.")
		return nil
	}

	fmt.Printf("Config file: %s\n", getConfigPath())
	fmt.Printf("Default profile: %s\n", cfg.DefaultProfile)
	fmt.Printf("Active profile: %s\n", resolveProfile(cmd))
	fmt.Println()

	table := NewTableWriter(os.Stdout, "PROFILE", "API KEY", "DEFAULT")
	for name, p := range cfg.Profiles {
		key := maskKey(p.APIKey)
		def := ""
		if name == cfg.DefaultProfile {
			def = "*"
		}
		table.AddRow(name, key, def)
	}
	table.Render()
	return nil
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func profileNames(cfg *CLIConfig) []string {
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	return names
}
