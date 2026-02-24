package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var domainsCmd = cli.Command{
	Name:  "domains",
	Usage: "Manage hosting domains",
	Commands: []*cli.Command{
		&domainsListCmd,
		&domainsGetCmd,
		&domainsAvailableCmd,
		&domainsContactsCmd,
		&domainsContactsUpdateCmd,
		&domainsContactUpdateCmd,
		&domainsContactVerifyCmd,
		&domainsDNSSECGetCmd,
		&domainsDNSSECUpdateCmd,
		&domainsHealthCmd,
		&domainsLocksGetCmd,
		&domainsLocksSetCmd,
		&domainsNameserversCmd,
		&domainsNameserversUpdateCmd,
		&domainsOverviewCmd,
		&domainsDNSCmd,
		&domainsDNSGetCmd,
		&domainsDNSCreateCmd,
		&domainsDNSUpdateAllCmd,
		&domainsDNSDeleteAllCmd,
		&domainsDNSUpdateCmd,
		&domainsDNSDeleteCmd,
		&domainsDNSExportCmd,
		&domainsDNSImportCmd,
		&domainsDNSValidateSetCmd,
		&domainsValidateZoneCmd,
		&domainsKeyRolloverCmd,
		&domainsDNSQueryMetricsCmd,
	},
	HideHelpCommand: true,
}

var domainsListCmd = cli.Command{
	Name:            "list",
	Usage:           "List domains",
	Flags:           PaginationFlags,
	Action:          handleDomainsList,
	HideHelpCommand: true,
}

func handleDomainsList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	domains := res.Get("domains")
	if !domains.Exists() || len(domains.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No domains found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "DOMAIN", "STATUS", "NAMESERVERS")
	domains.ForEach(func(_, d gjson.Result) bool {
		ns := ""
		d.Get("nameServers").ForEach(func(_, n gjson.Result) bool {
			if ns != "" {
				ns += ", "
			}
			ns += n.String()
			return true
		})
		table.AddRow(d.Get("domainName").String(), d.Get("status").String(), ns)
		return true
	})
	table.Render()
	return nil
}

var domainsGetCmd = cli.Command{
	Name:            "get",
	Usage:           "Get domain details",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsGet,
	HideHelpCommand: true,
}

func handleDomainsGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSCmd = cli.Command{
	Name:            "dns",
	Usage:           "List DNS records for a domain",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsDNS,
	HideHelpCommand: true,
}

func handleDomainsDNS(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets")
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	records := res.Get("resourceRecordSets")
	if !records.Exists() || len(records.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No DNS records found.")
		return nil
	}

	table := NewTableWriter(os.Stdout, "NAME", "TYPE", "TTL", "CONTENT")
	table.TruncOrder = []int{3, 0}
	records.ForEach(func(_, r gjson.Result) bool {
		content := ""
		r.Get("content").ForEach(func(_, c gjson.Result) bool {
			if content != "" {
				content += ", "
			}
			content += c.String()
			return true
		})
		table.AddRow(
			r.Get("name").String(),
			r.Get("type").String(),
			fmt.Sprintf("%d", r.Get("ttl").Int()),
			content,
		)
		return true
	})
	table.Render()
	return nil
}

var domainsDNSGetCmd = cli.Command{
	Name:            "dns-get",
	Usage:           "Get a specific DNS record set",
	ArgsUsage:       "<domain> <name> <type>",
	Action:          handleDomainsDNSGet,
	HideHelpCommand: true,
}

func handleDomainsDNSGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("domain, record name, and type required\nUsage: lw domains dns-get <domain> <name> <type>")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("/hosting/v2/domains/%s/resourceRecordSets/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSCreateCmd = cli.Command{
	Name:      "dns-create",
	Usage:     "Create a DNS record set",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "Record name", Required: true},
		&cli.StringFlag{Name: "type", Usage: "Record type (A, AAAA, CNAME, MX, TXT, NS, SRV)", Required: true},
		&cli.IntFlag{Name: "ttl", Usage: "TTL in seconds", Value: 3600},
		&cli.StringSliceFlag{Name: "content", Usage: "Record content (can be repeated)", Required: true},
	},
	Action:          handleDomainsDNSCreate,
	HideHelpCommand: true,
}

func handleDomainsDNSCreate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	payload := fmt.Sprintf(`{"name":%q,"type":%q,"ttl":%d,"content":%s}`,
		cmd.String("name"),
		cmd.String("type"),
		cmd.Int("ttl"),
		mustMarshalStrings(cmd.StringSlice("content")),
	)

	res, err := client.Post(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets", payload)
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSDeleteCmd = cli.Command{
	Name:            "dns-delete",
	Usage:           "Delete a DNS record set",
	ArgsUsage:       "<domain> <name> <type>",
	Action:          handleDomainsDNSDelete,
	HideHelpCommand: true,
}

func handleDomainsDNSDelete(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("domain, record name, and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("/hosting/v2/domains/%s/resourceRecordSets/%s/%s", args[0], args[1], args[2]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted DNS record %s %s for %s\n", args[1], args[2], args[0])
	return nil
}

var domainsAvailableCmd = cli.Command{
	Name:            "available",
	Usage:           "Check domain availability",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsAvailable,
	HideHelpCommand: true,
}

func handleDomainsAvailable(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/available")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsContactsCmd = cli.Command{
	Name:            "contacts",
	Usage:           "List domain contacts",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsContacts,
	HideHelpCommand: true,
}

func handleDomainsContacts(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/contacts")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsContactsUpdateCmd = cli.Command{
	Name:      "contacts-update",
	Usage:     "Update all contacts for a domain",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsContactsUpdate,
	HideHelpCommand: true,
}

func handleDomainsContactsUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/hosting/v2/domains/"+args[0]+"/contacts", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsContactUpdateCmd = cli.Command{
	Name:      "contact-update",
	Usage:     "Update a specific contact type for a domain",
	ArgsUsage: "<domain> <type>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsContactUpdate,
	HideHelpCommand: true,
}

func handleDomainsContactUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain name and contact type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/hosting/v2/domains/%s/contacts/%s", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsContactVerifyCmd = cli.Command{
	Name:      "contact-verify",
	Usage:     "Verify a contact for a domain",
	ArgsUsage: "<domain> <type>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsContactVerify,
	HideHelpCommand: true,
}

func handleDomainsContactVerify(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain name and contact type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/hosting/v2/domains/%s/contacts/%s/verify", args[0], args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSSECGetCmd = cli.Command{
	Name:            "dnssec",
	Usage:           "Inspect DNSSEC settings for a domain",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsDNSSECGet,
	HideHelpCommand: true,
}

func handleDomainsDNSSECGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/dnssec")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSSECUpdateCmd = cli.Command{
	Name:      "dnssec-update",
	Usage:     "Update DNSSEC settings for a domain",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsDNSSECUpdate,
	HideHelpCommand: true,
}

func handleDomainsDNSSECUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/hosting/v2/domains/"+args[0]+"/dnssec", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsHealthCmd = cli.Command{
	Name:            "health",
	Usage:           "Get domain health",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsHealth,
	HideHelpCommand: true,
}

func handleDomainsHealth(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/health")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsLocksGetCmd = cli.Command{
	Name:            "locks",
	Usage:           "Get domain locks",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsLocksGet,
	HideHelpCommand: true,
}

func handleDomainsLocksGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/locks")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsLocksSetCmd = cli.Command{
	Name:      "locks-set",
	Usage:     "Set domain locks",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsLocksSet,
	HideHelpCommand: true,
}

func handleDomainsLocksSet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/hosting/v2/domains/"+args[0]+"/locks", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsNameserversCmd = cli.Command{
	Name:            "nameservers",
	Usage:           "List nameservers for a domain",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsNameservers,
	HideHelpCommand: true,
}

func handleDomainsNameservers(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/nameservers")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsNameserversUpdateCmd = cli.Command{
	Name:      "nameservers-update",
	Usage:     "Update nameservers for a domain",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsNameserversUpdate,
	HideHelpCommand: true,
}

func handleDomainsNameserversUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/hosting/v2/domains/"+args[0]+"/nameservers", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsOverviewCmd = cli.Command{
	Name:            "overview",
	Usage:           "Get detailed domain information",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsOverview,
	HideHelpCommand: true,
}

func handleDomainsOverview(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/overview")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSUpdateAllCmd = cli.Command{
	Name:      "dns-update-all",
	Usage:     "Update all DNS records for a domain",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload with all records", Required: true},
	},
	Action:          handleDomainsDNSUpdateAll,
	HideHelpCommand: true,
}

func handleDomainsDNSUpdateAll(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSDeleteAllCmd = cli.Command{
	Name:            "dns-delete-all",
	Usage:           "Delete all DNS records for a domain",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsDNSDeleteAll,
	HideHelpCommand: true,
}

func handleDomainsDNSDeleteAll(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted all DNS records for %s\n", args[0])
	return nil
}

var domainsDNSUpdateCmd = cli.Command{
	Name:      "dns-update",
	Usage:     "Update a specific DNS record",
	ArgsUsage: "<domain> <name> <type>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsDNSUpdate,
	HideHelpCommand: true,
}

func handleDomainsDNSUpdate(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 3 {
		return fmt.Errorf("domain, record name, and type required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("/hosting/v2/domains/%s/resourceRecordSets/%s/%s", args[0], args[1], args[2]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSExportCmd = cli.Command{
	Name:            "dns-export",
	Usage:           "Export DNS records as bind file content",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsDNSExport,
	HideHelpCommand: true,
}

func handleDomainsDNSExport(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets/export")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSImportCmd = cli.Command{
	Name:      "dns-import",
	Usage:     "Import DNS records from bind file content",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload with bind content", Required: true},
	},
	Action:          handleDomainsDNSImport,
	HideHelpCommand: true,
}

func handleDomainsDNSImport(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets/import", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSValidateSetCmd = cli.Command{
	Name:      "dns-validate-set",
	Usage:     "Validate a resource record set",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsDNSValidateSet,
	HideHelpCommand: true,
}

func handleDomainsDNSValidateSet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/hosting/v2/domains/"+args[0]+"/resourceRecordSets/validateSet", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsValidateZoneCmd = cli.Command{
	Name:      "validate-zone",
	Usage:     "Validate zone for a domain",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsValidateZone,
	HideHelpCommand: true,
}

func handleDomainsValidateZone(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/hosting/v2/domains/"+args[0]+"/validateZone", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsKeyRolloverCmd = cli.Command{
	Name:      "key-rollover",
	Usage:     "Perform key rollover for a domain",
	ArgsUsage: "<domain>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "payload", Usage: "JSON payload", Required: true},
	},
	Action:          handleDomainsKeyRollover,
	HideHelpCommand: true,
}

func handleDomainsKeyRollover(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/hosting/v2/domains/"+args[0]+"/keyRollover", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

var domainsDNSQueryMetricsCmd = cli.Command{
	Name:            "dns-query-metrics",
	Usage:           "Show DNS query metrics for a domain",
	ArgsUsage:       "<domain>",
	Action:          handleDomainsDNSQueryMetrics,
	HideHelpCommand: true,
}

func handleDomainsDNSQueryMetrics(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain name required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/hosting/v2/domains/"+args[0]+"/metrics/dnsQuery")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("format"), cmd.Root().String("transform"))
}

func mustMarshalStrings(s []string) string {
	b, _ := gjson.Parse("[]").Value().([]interface{})
	_ = b
	out := "["
	for i, v := range s {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf("%q", v)
	}
	out += "]"
	return out
}
