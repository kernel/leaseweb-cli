package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var emailsCmd = cli.Command{
	Name:    "emails",
	Aliases: []string{"email"},
	Usage:   "Manage email services",
	Commands: []*cli.Command{
		&emailDomainsListCmd, &emailDomainCreateCmd, &emailDomainGetCmd, &emailDomainUpdateCmd, &emailDomainDeleteCmd,
		&emailDomainVerifyCmd, &emailDomainDNSCmd,
		&emailMailboxesListCmd, &emailMailboxCreateCmd, &emailMailboxGetCmd, &emailMailboxUpdateCmd, &emailMailboxDeleteCmd,
		&emailMailboxAutoReplyCmd, &emailMailboxAutoReplyUpdateCmd,
		&emailForwardsListCmd, &emailForwardCreateCmd, &emailForwardGetCmd, &emailForwardUpdateCmd, &emailForwardDeleteCmd,
		&emailAliasesListCmd, &emailAliasCreateCmd, &emailAliasGetCmd, &emailAliasDeleteCmd,
		&emailSpamFilterCmd, &emailSpamFilterUpdateCmd,
	},
	HideHelpCommand: true,
}

func emailDomainPath(args []string) string { return "/email/v2/domains/" + args[0] }

// Domains
var emailDomainsListCmd = cli.Command{Name: "domains", Usage: "List email domains", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/email/v2/domains?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailDomainCreateCmd = cli.Command{Name: "domain-create", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, "/email/v2/domains", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailDomainGetCmd = cli.Command{Name: "domain-get", ArgsUsage: "<domain>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailDomainPath(args))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailDomainUpdateCmd = cli.Command{Name: "domain-update", ArgsUsage: "<domain>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, emailDomainPath(args), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailDomainDeleteCmd = cli.Command{Name: "domain-delete", ArgsUsage: "<domain>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, emailDomainPath(args))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted domain %s\n", args[0])
	return nil
}, HideHelpCommand: true}

var emailDomainVerifyCmd = cli.Command{Name: "domain-verify", ArgsUsage: "<domain>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Post(ctx, emailDomainPath(args)+"/verify", "")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailDomainDNSCmd = cli.Command{Name: "domain-dns", ArgsUsage: "<domain>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailDomainPath(args)+"/dnsRecords")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

// Mailboxes
func emailMailboxPath(args []string) string {
	return fmt.Sprintf("/email/v2/domains/%s/mailboxes/%s", args[0], args[1])
}

var emailMailboxesListCmd = cli.Command{Name: "mailboxes", ArgsUsage: "<domain>", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailDomainPath(args)+"/mailboxes?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailMailboxCreateCmd = cli.Command{Name: "mailbox-create", ArgsUsage: "<domain>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, emailDomainPath(args)+"/mailboxes", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailMailboxGetCmd = cli.Command{Name: "mailbox-get", ArgsUsage: "<domain> <mailbox>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and mailbox required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailMailboxPath(args))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailMailboxUpdateCmd = cli.Command{Name: "mailbox-update", ArgsUsage: "<domain> <mailbox>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and mailbox required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, emailMailboxPath(args), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailMailboxDeleteCmd = cli.Command{Name: "mailbox-delete", ArgsUsage: "<domain> <mailbox>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and mailbox required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, emailMailboxPath(args))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted mailbox %s\n", args[1])
	return nil
}, HideHelpCommand: true}

var emailMailboxAutoReplyCmd = cli.Command{Name: "mailbox-autoreply", ArgsUsage: "<domain> <mailbox>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and mailbox required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailMailboxPath(args)+"/autoReply")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailMailboxAutoReplyUpdateCmd = cli.Command{Name: "mailbox-autoreply-update", ArgsUsage: "<domain> <mailbox>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and mailbox required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, emailMailboxPath(args)+"/autoReply", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

// Forwards
var emailForwardsListCmd = cli.Command{Name: "forwards", ArgsUsage: "<domain>", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailDomainPath(args)+"/forwards?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailForwardCreateCmd = cli.Command{Name: "forward-create", ArgsUsage: "<domain>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, emailDomainPath(args)+"/forwards", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailForwardGetCmd = cli.Command{Name: "forward-get", ArgsUsage: "<domain> <forward-id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and forward ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/forwards/%s", emailDomainPath(args), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailForwardUpdateCmd = cli.Command{Name: "forward-update", ArgsUsage: "<domain> <forward-id>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and forward ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, fmt.Sprintf("%s/forwards/%s", emailDomainPath(args), args[1]), []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailForwardDeleteCmd = cli.Command{Name: "forward-delete", ArgsUsage: "<domain> <forward-id>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and forward ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("%s/forwards/%s", emailDomainPath(args), args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted forward %s\n", args[1])
	return nil
}, HideHelpCommand: true}

// Aliases
var emailAliasesListCmd = cli.Command{Name: "aliases", ArgsUsage: "<domain>", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailDomainPath(args)+"/aliases?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailAliasCreateCmd = cli.Command{Name: "alias-create", ArgsUsage: "<domain>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PostJSON(ctx, emailDomainPath(args)+"/aliases", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailAliasGetCmd = cli.Command{Name: "alias-get", ArgsUsage: "<domain> <alias>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and alias required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, fmt.Sprintf("%s/aliases/%s", emailDomainPath(args), args[1]))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailAliasDeleteCmd = cli.Command{Name: "alias-delete", ArgsUsage: "<domain> <alias>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 2 {
		return fmt.Errorf("domain and alias required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	_, err = client.Delete(ctx, fmt.Sprintf("%s/aliases/%s", emailDomainPath(args), args[1]))
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Deleted alias %s\n", args[1])
	return nil
}, HideHelpCommand: true}

// Spam Filter
var emailSpamFilterCmd = cli.Command{Name: "spam-filter", ArgsUsage: "<domain>", Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, emailDomainPath(args)+"/spamFilter")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var emailSpamFilterUpdateCmd = cli.Command{Name: "spam-filter-update", ArgsUsage: "<domain>", Flags: payloadFlag, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("domain required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.PutJSON(ctx, emailDomainPath(args)+"/spamFilter", []byte(cmd.String("payload")))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}
