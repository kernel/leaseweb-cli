# lw — CLI for the Leaseweb API

A command-line interface for managing [Leaseweb](https://www.leaseweb.com/) infrastructure: dedicated servers, public cloud instances, VPS, virtual servers, private clouds, CDN, DNS, email, and more.

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

For managing multiple Leaseweb accounts, use profiles. Run `lw config init` to set up a profile:

```sh
$ lw config init
Profile name, e.g. "us" or "ca": us
API key for "us": 74B196B1-...
Set "us" as the default profile? (y/n) [y]: y
Profile "us" saved to ~/.config/lw/config.yaml

$ lw config init
Existing profiles: us

Profile name, e.g. "us" or "ca": ca
API key for "ca": BD483105-...
Set "ca" as the default profile? (y/n) [n]: n
Profile "ca" saved to ~/.config/lw/config.yaml
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

**Resolution order** for API key: `LEASEWEB_API_KEY` env var > profile config.

## Usage

```
NAME:
   lw - CLI for the Leaseweb API

USAGE:
   lw [global options] [command [command options]]

VERSION:
   dev

COMMANDS:
   abuse-reports, abuse    Manage abuse reports
   acronis-backup, backup  Manage Acronis backup
   aggregation-packs, ap   Manage aggregation packs
   api-keys, keys          Manage API keys
   cdn, c                  Manage CDN resources
   colocations, colo       Manage colocations
   config                  Manage CLI configuration
   datacenter-access, dca  Manage datacenter access requests
   dedicated-racks, dr     Manage dedicated racks
   dedicated-servers, ds   Manage dedicated servers
   domains                 Manage hosting domains
   emails, email           Manage email services
   floating-ips, fip       Manage floating IPs
   instances, i            Manage public cloud instances
   invoices                Manage invoices
   ips                     Manage IP addresses
   load-balancers, lb      Manage public cloud load balancers
   network-equipment, ne   Manage dedicated network equipment
   private-clouds, pc      Manage private clouds
   private-networks, pn    Manage private networks
   remote-management, rm   Manage OpenVPN remote management
   services                Manage services
   storage                 Manage storage
   traffic-policy, tp      Manage traffic policies
   virtual-servers, vs     Manage virtual servers
   vps, v                  Manage VPS instances
   webhosting, wh          Manage webhosting packages

GLOBAL OPTIONS:
   --profile string, -p string  Config profile to use
   --debug                      Enable debug logging of HTTP requests
   --output string, -o string   Output format (one of: auto, json, jsonline, pretty, raw, yaml) (default: "auto")
   --transform string           GJSON expression to transform output
   --help, -h                   show help
   --version, -v                print the version
```

### Examples

```sh
# List dedicated servers
lw dedicated-servers list

# Get server details as YAML
lw ds get 12490707 -o yaml

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

| Format     | Description                          |
|------------|--------------------------------------|
| `auto`     | Table for lists, structured text for details (default) |
| `json`     | Pretty-printed JSON with syntax colors |
| `jsonline` | Compact JSON, single line (useful for piping) |
| `pretty`   | Pretty-printed JSON without colors   |
| `raw`      | Raw JSON as returned by the API      |
| `yaml`     | YAML output                          |

### Transform

The `--transform` flag accepts [GJSON](https://github.com/tidwall/gjson) expressions to extract or query nested data:

```sh
# Get just the CPU info
lw ds get 12490707 --transform "specs.cpu"

# Get all server IDs from a list
lw ds list -o raw --transform "servers.#.id"
```

## Subcommands

Use `lw <command> --help` for details on any subcommand. Here's a summary:

| Command | Alias | Description |
|---------|-------|-------------|
| `abuse-reports` | `abuse` | List, get, resolve abuse reports, manage messages and attachments |
| `acronis-backup` | `backup` | List backup items, get details, view metrics |
| `aggregation-packs` | `ap` | List and get aggregation packs |
| `api-keys` | `keys` | CRUD API keys, validate keys, list capabilities |
| `cdn` | `c` | Distributions, origins, cache, SSL, WAF, geo-restrictions, metrics |
| `colocations` | `colo` | CRUD colocations, credentials, IPs, metrics, notifications |
| `config` | | init, show |
| `datacenter-access` | `dca` | Access requests, datacenters, contacts, visitors |
| `dedicated-racks` | `dr` | CRUD racks, credentials, IPs, metrics, notifications |
| `dedicated-servers` | `ds` | Full server lifecycle, credentials, IPs, jobs, metrics, DHCP, notifications |
| `domains` | | DNS records, DNSSEC, nameservers, contacts, locks, zone import/export |
| `emails` | `email` | Domains, mailboxes, forwards, aliases, spam filter, auto-reply |
| `floating-ips` | `fip` | CRUD floating IP ranges, definitions, assign/unassign |
| `instances` | `i` | Full instance lifecycle, credentials, IPs, snapshots, ISOs, security groups |
| `invoices` | | List, get, PDF download, proforma, CSV export |
| `ips` | | List, get, update, null route, reverse lookup (IPv4 + IPv6) |
| `load-balancers` | `lb` | CRUD, listeners, IPs, metrics, monitoring |
| `network-equipment` | `ne` | CRUD equipment, credentials, IPs, power, null routes |
| `private-clouds` | `pc` | CRUD private clouds, credentials, metrics |
| `private-networks` | `pn` | CRUD networks, servers, DHCP reservations |
| `remote-management` | `rm` | OpenVPN profiles, credentials |
| `services` | | List, get, update, cancel/uncancel |
| `storage` | | List storage, VMs, volumes, grow volumes |
| `traffic-policy` | `tp` | List, get, update policies, history, reset |
| `virtual-servers` | `vs` | CRUD servers, credentials, metrics, snapshots, templates |
| `vps` | `v` | Full VPS lifecycle, credentials, IPs, snapshots, monitoring, notifications |
| `webhosting` | `wh` | Packages, usernames, domain aliases, catch-all |

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
