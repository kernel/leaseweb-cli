package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v3"
)

type ProfileConfig struct {
	APIKey string `koanf:"api_key"`
}

type CLIConfig struct {
	DefaultProfile string                   `koanf:"default_profile"`
	Profiles       map[string]ProfileConfig `koanf:"profiles"`
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "lw")
}

func getConfigPath() string {
	dir := getConfigDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "config.yaml")
}

func loadConfig() *CLIConfig {
	cfg := &CLIConfig{
		Profiles: make(map[string]ProfileConfig),
	}
	k := koanf.New(".")

	configPath := getConfigPath()
	if configPath != "" {
		_ = k.Load(file.Provider(configPath), yaml.Parser())
	}
	_ = k.Unmarshal("", cfg)
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]ProfileConfig)
	}
	return cfg
}

func resolveProfile(cmd *cli.Command) (string, error) {
	if p := cmd.Root().String("profile"); p != "" {
		return p, nil
	}
	if p := os.Getenv("LEASEWEB_PROFILE"); p != "" {
		return p, nil
	}
	cfg := loadConfig()
	if cfg.DefaultProfile != "" {
		return cfg.DefaultProfile, nil
	}
	return "", fmt.Errorf("no profile specified and no default profile set. Use -p <profile>, set LEASEWEB_PROFILE, or run 'lw config init'")
}

func resolveAPIKey(cmd *cli.Command) (string, error) {
	if k := os.Getenv("LEASEWEB_API_KEY"); k != "" {
		return k, nil
	}
	profile, err := resolveProfile(cmd)
	if err != nil {
		return "", err
	}
	cfg := loadConfig()
	if p, ok := cfg.Profiles[profile]; ok && p.APIKey != "" {
		return p.APIKey, nil
	}
	return "", fmt.Errorf("no API key found for profile %q. Set LEASEWEB_API_KEY or run 'lw config init'", profile)
}

func resolveBaseURL() string {
	if u := os.Getenv("LEASEWEB_BASE_URL"); u != "" {
		return u
	}
	return "https://api.leaseweb.com"
}

func writeConfig(cfg *CLIConfig) error {
	dir := getConfigDir()
	if dir == "" {
		return fmt.Errorf("could not determine home directory")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	var b strings.Builder
	if cfg.DefaultProfile != "" {
		fmt.Fprintf(&b, "default_profile: %s\n", cfg.DefaultProfile)
	}
	if len(cfg.Profiles) > 0 {
		b.WriteString("profiles:\n")
		for name, p := range cfg.Profiles {
			fmt.Fprintf(&b, "  %s:\n", name)
			fmt.Fprintf(&b, "    api_key: %q\n", p.APIKey)
		}
	}

	path := getConfigPath()
	return os.WriteFile(path, []byte(b.String()), 0600)
}
