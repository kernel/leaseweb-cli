# lw — CLI for the Leaseweb API

A command-line interface for managing [Leaseweb](https://www.leaseweb.com/) infrastructure: dedicated servers, public cloud instances, IP addresses, invoices, DNS, and more.

## Installation

### Homebrew

```sh
brew install kernel/tap/lw
```

### Go

```sh
go install github.com/kernel/leaseweb-cli/cmd/lw@latest
```

### Binary releases

Download pre-built binaries from the [Releases](https://github.com/kernel/leaseweb-cli/releases) page.

## Configuration

### Quick start

Set your API key as an environment variable:

```sh
export LEASEWEB_API_KEY="your-api-key-here"
lw dedicated-servers list
```

### Profiles

For managing multiple Leaseweb accounts, use profiles. Run `lw config init` to set them up interactively:

```sh
$ lw config init
Leaseweb CLI Configuration

Profile name (or 'done' to finish) [default]: us
API key for "us": 74B196B1-...
Profile "us" saved.
Profile name (or 'done' to finish) [default]: ca
API key for "ca": BD483105-...
Profile "ca" saved.
Profile name (or 'done' to finish) [default]: done
Default profile [us]: us

Configuration written to ~/.config/lw/config.yaml
```

Then switch profiles with `-p`:

```sh
lw -p us dedicated-servers list
lw -p ca dedicated-servers list
```

The config file lives at `~/.config/lw/config.yaml`:

```yaml
default_profile: us
profiles:
  us:
    api_key: "74B196B1-..."
  ca:
    api_key: "BD483105-..."
```

**Resolution order** for API key: `--api-key` flag > `LEASEWEB_API_KEY` env > profile config.

## Usage

```
NAME:
   lw - CLI for the Leaseweb API

USAGE:
   lw [global options] [command [command options]]

COMMANDS:
   config                 Manage CLI configuration
   dedicated-servers, ds  Manage dedicated servers
   domains                Manage hosting domains
   instances, i           Manage public cloud instances
   invoices               Manage invoices
   ips                    Manage IP addresses
   load-balancers, lb     Manage public cloud load balancers
   private-networks, pn   Manage private networks
   services               Manage services

GLOBAL OPTIONS:
   --profile string, -p string  Config profile to use
   --api-key string             Leaseweb API key (overrides profile)
   --base-url string            Override the base URL for API requests
   --debug                      Enable debug logging of HTTP requests
   --format string              Output format: auto, json, pretty, raw, yaml (default: "auto")
   --transform string           GJSON expression to transform output
   --help, -h                   show help
   --version, -v                print the version
```

### Examples

```sh
# List dedicated servers
lw dedicated-servers list

# Get server details as YAML
lw ds get 12490707 --format yaml

# Extract a specific field with --transform (uses GJSON syntax)
lw ds get 12490707 --transform "specs.cpu"

# List invoices
lw invoices list

# Download an invoice PDF
lw invoices pdf 84048268 -o invoice.pdf

# List IPs, filtering by version
lw ips list --version 4

# Null route an IP
lw ips null-route 1.2.3.4 --comment "under attack"

# List services
lw services list

# Public cloud instances
lw instances list
lw instances regions
lw instances types --region eu-west-3

# DNS management
lw domains list
lw domains dns example.com
lw domains dns-create example.com --name www --type A --content 1.2.3.4

# Debug HTTP requests
lw --debug ds list
```

### Output formats

| Format   | Description                          |
|----------|--------------------------------------|
| `auto`   | Table for list commands, JSON for detail commands (default) |
| `json`   | Pretty-printed JSON with colors      |
| `pretty` | Pretty-printed JSON without colors   |
| `raw`    | Compact JSON, one line               |
| `yaml`   | YAML output                          |

### Transform

The `--transform` flag accepts [GJSON](https://github.com/tidwall/gjson) expressions to extract or query nested data:

```sh
# Get just the CPU info
lw ds get 12490707 --transform "specs.cpu"

# Get all server IDs from a list
lw ds list --format raw --transform "servers.#.id"
```

## Subcommands

Use `lw <command> --help` for details on any subcommand. Here's a summary:

| Command | Alias | Subcommands |
|---------|-------|-------------|
| `dedicated-servers` | `ds` | list, get, update, power-on/off/cycle, power-status, rescue, install, credentials, jobs, hardware-info, metrics, network-interfaces |
| `instances` | `i` | list, get, launch, terminate, start, stop, reboot, update, console, credentials, ips, snapshots, metrics, regions, types, images |
| `ips` | | list, get, update, null-route, remove-null-route, null-route-history, reverse-lookup |
| `invoices` | | list, get, pdf, proforma |
| `services` | | list, get, update, cancel, uncancel, cancellation-reasons |
| `domains` | | list, get, dns, dns-get, dns-create, dns-delete |
| `load-balancers` | `lb` | list, get, create, update, delete, listeners |
| `private-networks` | `pn` | list, get, create, update, delete, servers |
| `config` | | init, show |

## Development

```sh
# Build
make build

# Install to $GOPATH/bin
make install

# Run tests
make test

# Lint
make lint

# Vet
make vet
```

## License

Apache-2.0
