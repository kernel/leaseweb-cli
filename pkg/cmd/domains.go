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
		&domainsDNSCmd,
		&domainsDNSGetCmd,
		&domainsDNSCreateCmd,
		&domainsDNSDeleteCmd,
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
