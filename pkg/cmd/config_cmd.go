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

	if len(cfg.Profiles) > 0 {
		fmt.Printf("Existing profiles: %s\n", strings.Join(profileNames(cfg), ", "))
		fmt.Println()
	}

	name := prompt(reader, `Profile name, e.g. "us" or "ca"`, "")
	if name == "" {
		return fmt.Errorf("profile name is required")
	}

	if p, ok := cfg.Profiles[name]; ok {
		fmt.Printf("Profile %q already exists (API key: %s)\n", name, maskKey(p.APIKey))
		overwrite := prompt(reader, "Overwrite? (y/n)", "n")
		if strings.ToLower(overwrite) != "y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	apiKey := prompt(reader, fmt.Sprintf("API key for %q", name), "")
	if apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	cfg.Profiles[name] = ProfileConfig{APIKey: apiKey}

	defaultDefault := "n"
	if cfg.DefaultProfile == "" || len(cfg.Profiles) == 1 {
		defaultDefault = "y"
	}
	setDefault := prompt(reader, fmt.Sprintf("Set %q as the default profile? (y/n)", name), defaultDefault)
	if strings.ToLower(setDefault) == "y" {
		cfg.DefaultProfile = name
	}

	if err := writeConfig(cfg); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Printf("Profile %q saved to %s\n", name, getConfigPath())
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
	active := cfg.DefaultProfile
	if p, err := resolveProfile(cmd); err == nil {
		active = p
	}
	fmt.Printf("Active profile: %s\n", active)
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
