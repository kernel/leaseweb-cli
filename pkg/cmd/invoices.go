package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/urfave/cli/v3"
)

var invoicesCmd = cli.Command{
	Name:  "invoices",
	Usage: "Manage invoices",
	Commands: []*cli.Command{
		&invoicesListCmd,
		&invoicesGetCmd,
		&invoicesPDFCmd,
		&invoicesExportCSVCmd,
		&invoicesProformaCmd,
	},
	HideHelpCommand: true,
}

var invoicesListCmd = cli.Command{
	Name:  "list",
	Usage: "List invoices",
	Flags: PaginationFlags,
	Action:          handleInvoicesList,
	HideHelpCommand: true,
}

func handleInvoicesList(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/invoices/v1/invoices?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}

	format := cmd.Root().String("output")
	if format != "auto" {
		return ShowResult(os.Stdout, res, format, cmd.Root().String("transform"))
	}

	invoices := res.Get("invoices")
	if !invoices.Exists() || len(invoices.Array()) == 0 {
		fmt.Fprintln(os.Stderr, "No invoices found.")
		return nil
	}

	items := invoices.Array()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Get("date").String() > items[j].Get("date").String()
	})

	table := NewTableWriter(os.Stdout, "ID", "DATE", "STATUS", "TOTAL", "CURRENCY", "DUE DATE")
	for _, inv := range items {
		table.AddRow(
			inv.Get("id").String(),
			dateOnly(inv.Get("date").String()),
			inv.Get("status").String(),
			fmt.Sprintf("%.2f", inv.Get("total").Float()),
			inv.Get("currency").String(),
			dateOnly(inv.Get("dueDate").String()),
		)
	}
	table.Render()
	return nil
}

func dateOnly(s string) string {
	if i := strings.IndexByte(s, 'T'); i > 0 {
		return s[:i]
	}
	return s
}

var invoicesGetCmd = cli.Command{
	Name:      "get",
	Usage:     "Get invoice details",
	ArgsUsage: "<invoice-id>",
	Action:    handleInvoicesGet,
	HideHelpCommand: true,
}

func handleInvoicesGet(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("invoice ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/invoices/v1/invoices/"+args[0])
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}

var invoicesPDFCmd = cli.Command{
	Name:      "pdf",
	Usage:     "Download invoice PDF",
	ArgsUsage: "<invoice-id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output file path (default: <invoice-id>.pdf)",
		},
	},
	Action:          handleInvoicesPDF,
	HideHelpCommand: true,
}

func handleInvoicesPDF(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("invoice ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}

	data, _, err := client.DoRaw(ctx, "GET", "/invoices/v1/invoices/"+args[0]+"/pdf")
	if err != nil {
		return err
	}

	output := cmd.String("output")
	if output == "" {
		output = args[0] + ".pdf"
	}

	if err := os.WriteFile(output, data, 0644); err != nil {
		return fmt.Errorf("writing PDF: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Downloaded invoice to %s\n", output)
	return nil
}

var invoicesExportCSVCmd = cli.Command{
	Name:  "export-csv",
	Usage: "Export invoices as CSV",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output file path (default: invoices.csv)",
		},
	},
	Action:          handleInvoicesExportCSV,
	HideHelpCommand: true,
}

func handleInvoicesExportCSV(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	data, _, err := client.DoRaw(ctx, "GET", "/invoices/v1/invoices/export/csv")
	if err != nil {
		return err
	}
	output := cmd.String("output")
	if output == "" {
		output = "invoices.csv"
	}
	if err := os.WriteFile(output, data, 0644); err != nil {
		return fmt.Errorf("writing CSV: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Exported invoices to %s\n", output)
	return nil
}

var invoicesProformaCmd = cli.Command{
	Name:            "proforma",
	Usage:           "Get pro forma invoice",
	Action:          handleInvoicesProforma,
	HideHelpCommand: true,
}

func handleInvoicesProforma(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/invoices/v1/proforma")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}
