package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var cdnCmd = cli.Command{
	Name:    "cdn",
	Aliases: []string{"c"},
	Usage:   "Manage CDN resources",
	Commands: []*cli.Command{
		&cdnDistributionsListCmd, &cdnDistributionCreateCmd, &cdnDistributionGetCmd, &cdnDistributionUpdateCmd, &cdnDistributionDeleteCmd,
		&cdnOriginsCmd, &cdnOriginCreateCmd, &cdnOriginGetCmd, &cdnOriginUpdateCmd, &cdnOriginDeleteCmd,
		&cdnCacheSettingsCmd, &cdnCacheSettingCreateCmd, &cdnCacheSettingGetCmd, &cdnCacheSettingUpdateCmd, &cdnCacheSettingDeleteCmd,
		&cdnCachePurgeCmd,
		&cdnSSLCertsCmd, &cdnSSLCertCreateCmd, &cdnSSLCertGetCmd, &cdnSSLCertUpdateCmd, &cdnSSLCertDeleteCmd,
		&cdnCustomRulesCmd, &cdnCustomRuleCreateCmd, &cdnCustomRuleGetCmd, &cdnCustomRuleUpdateCmd, &cdnCustomRuleDeleteCmd,
		&cdnWAFCmd, &cdnWAFUpdateCmd, &cdnWAFIPsCmd, &cdnWAFIPCreateCmd, &cdnWAFIPDeleteCmd,
		&cdnGeoRestrictionsCmd, &cdnGeoRestrictionsUpdateCmd,
		&cdnRateLimitCmd, &cdnRateLimitUpdateCmd,
		&cdnMetricsBandwidthCmd, &cdnMetricsRequestsCmd, &cdnMetricsStatusCodesCmd,
		&cdnMetricsCacheHitCmd, &cdnMetricsEdgeLocationsCmd, &cdnMetricsOriginBandwidthCmd,
		&cdnAccessLogsCmd,
		&cdnEdgeLocationsCmd,
		&cdnTokenAuthCmd, &cdnTokenAuthUpdateCmd,
		&cdnHotlinkProtectionCmd, &cdnHotlinkProtectionUpdateCmd,
	},
	HideHelpCommand: true,
}

func cdnCRUD(basePath string, minArgs int, argErr string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < minArgs {
			return fmt.Errorf("%s", argErr)
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		path := basePath
		for i := 0; i < minArgs; i++ {
			path = fmt.Sprintf(path, args[i])
		}
		res, err := client.Get(ctx, path)
		if err != nil {
			return err
		}
		return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
	}
}

func cdnPost(basePath string, minArgs int, argErr string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < minArgs {
			return fmt.Errorf("%s", argErr)
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		path := basePath
		for i := 0; i < minArgs; i++ {
			path = fmt.Sprintf(path, args[i])
		}
		res, err := client.PostJSON(ctx, path, []byte(cmd.String("payload")))
		if err != nil {
			return err
		}
		return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
	}
}

func cdnPut(basePath string, minArgs int, argErr string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < minArgs {
			return fmt.Errorf("%s", argErr)
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		path := basePath
		for i := 0; i < minArgs; i++ {
			path = fmt.Sprintf(path, args[i])
		}
		res, err := client.PutJSON(ctx, path, []byte(cmd.String("payload")))
		if err != nil {
			return err
		}
		return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
	}
}

func cdnDel(basePath string, minArgs int, argErr string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < minArgs {
			return fmt.Errorf("%s", argErr)
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		path := basePath
		for i := 0; i < minArgs; i++ {
			path = fmt.Sprintf(path, args[i])
		}
		_, err = client.Delete(ctx, path)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Deleted\n")
		return nil
	}
}

// Distributions
var cdnDistributionsListCmd = cli.Command{Name: "list", Usage: "List distributions", Flags: PaginationFlags, Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cdn/v2/distributions?"+PaginationQuery(cmd))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

var cdnDistributionCreateCmd = cli.Command{Name: "create", Usage: "Create distribution", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions", 0, ""), HideHelpCommand: true}
var cdnDistributionGetCmd = cli.Command{Name: "get", ArgsUsage: "<id>", Action: cdnCRUD("/cdn/v2/distributions/%s", 1, "distribution ID required"), HideHelpCommand: true}
var cdnDistributionUpdateCmd = cli.Command{Name: "update", ArgsUsage: "<id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s", 1, "distribution ID required"), HideHelpCommand: true}
var cdnDistributionDeleteCmd = cli.Command{Name: "delete", ArgsUsage: "<id>", Action: cdnDel("/cdn/v2/distributions/%s", 1, "distribution ID required"), HideHelpCommand: true}

// Origins
var cdnOriginsCmd = cli.Command{Name: "origins", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/origins", 1, "distribution ID required"), HideHelpCommand: true}
var cdnOriginCreateCmd = cli.Command{Name: "origin-create", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions/%s/origins", 1, "distribution ID required"), HideHelpCommand: true}
var cdnOriginGetCmd = cli.Command{Name: "origin-get", ArgsUsage: "<dist-id> <origin-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/origins/%s", 2, "distribution ID and origin ID required"), HideHelpCommand: true}
var cdnOriginUpdateCmd = cli.Command{Name: "origin-update", ArgsUsage: "<dist-id> <origin-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/origins/%s", 2, "distribution ID and origin ID required"), HideHelpCommand: true}
var cdnOriginDeleteCmd = cli.Command{Name: "origin-delete", ArgsUsage: "<dist-id> <origin-id>", Action: cdnDel("/cdn/v2/distributions/%s/origins/%s", 2, "distribution ID and origin ID required"), HideHelpCommand: true}

// Cache Settings
var cdnCacheSettingsCmd = cli.Command{Name: "cache-settings", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/cacheSettings", 1, "distribution ID required"), HideHelpCommand: true}
var cdnCacheSettingCreateCmd = cli.Command{Name: "cache-setting-create", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions/%s/cacheSettings", 1, "distribution ID required"), HideHelpCommand: true}
var cdnCacheSettingGetCmd = cli.Command{Name: "cache-setting-get", ArgsUsage: "<dist-id> <cs-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/cacheSettings/%s", 2, "distribution ID and setting ID required"), HideHelpCommand: true}
var cdnCacheSettingUpdateCmd = cli.Command{Name: "cache-setting-update", ArgsUsage: "<dist-id> <cs-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/cacheSettings/%s", 2, "distribution ID and setting ID required"), HideHelpCommand: true}
var cdnCacheSettingDeleteCmd = cli.Command{Name: "cache-setting-delete", ArgsUsage: "<dist-id> <cs-id>", Action: cdnDel("/cdn/v2/distributions/%s/cacheSettings/%s", 2, "distribution ID and setting ID required"), HideHelpCommand: true}

// Cache Purge
var cdnCachePurgeCmd = cli.Command{Name: "cache-purge", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions/%s/cachePurge", 1, "distribution ID required"), HideHelpCommand: true}

// SSL Certificates
var cdnSSLCertsCmd = cli.Command{Name: "ssl-certs", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/sslCertificates", 1, "distribution ID required"), HideHelpCommand: true}
var cdnSSLCertCreateCmd = cli.Command{Name: "ssl-cert-create", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions/%s/sslCertificates", 1, "distribution ID required"), HideHelpCommand: true}
var cdnSSLCertGetCmd = cli.Command{Name: "ssl-cert-get", ArgsUsage: "<dist-id> <cert-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/sslCertificates/%s", 2, "distribution ID and cert ID required"), HideHelpCommand: true}
var cdnSSLCertUpdateCmd = cli.Command{Name: "ssl-cert-update", ArgsUsage: "<dist-id> <cert-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/sslCertificates/%s", 2, "distribution ID and cert ID required"), HideHelpCommand: true}
var cdnSSLCertDeleteCmd = cli.Command{Name: "ssl-cert-delete", ArgsUsage: "<dist-id> <cert-id>", Action: cdnDel("/cdn/v2/distributions/%s/sslCertificates/%s", 2, "distribution ID and cert ID required"), HideHelpCommand: true}

// Custom Rules
var cdnCustomRulesCmd = cli.Command{Name: "custom-rules", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/customRules", 1, "distribution ID required"), HideHelpCommand: true}
var cdnCustomRuleCreateCmd = cli.Command{Name: "custom-rule-create", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions/%s/customRules", 1, "distribution ID required"), HideHelpCommand: true}
var cdnCustomRuleGetCmd = cli.Command{Name: "custom-rule-get", ArgsUsage: "<dist-id> <rule-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/customRules/%s", 2, "distribution ID and rule ID required"), HideHelpCommand: true}
var cdnCustomRuleUpdateCmd = cli.Command{Name: "custom-rule-update", ArgsUsage: "<dist-id> <rule-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/customRules/%s", 2, "distribution ID and rule ID required"), HideHelpCommand: true}
var cdnCustomRuleDeleteCmd = cli.Command{Name: "custom-rule-delete", ArgsUsage: "<dist-id> <rule-id>", Action: cdnDel("/cdn/v2/distributions/%s/customRules/%s", 2, "distribution ID and rule ID required"), HideHelpCommand: true}

// WAF
var cdnWAFCmd = cli.Command{Name: "waf", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/waf", 1, "distribution ID required"), HideHelpCommand: true}
var cdnWAFUpdateCmd = cli.Command{Name: "waf-update", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/waf", 1, "distribution ID required"), HideHelpCommand: true}
var cdnWAFIPsCmd = cli.Command{Name: "waf-ips", Usage: "List WAF IP whitelist", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/waf/ipWhitelist", 1, "distribution ID required"), HideHelpCommand: true}
var cdnWAFIPCreateCmd = cli.Command{Name: "waf-ip-create", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPost("/cdn/v2/distributions/%s/waf/ipWhitelist", 1, "distribution ID required"), HideHelpCommand: true}
var cdnWAFIPDeleteCmd = cli.Command{Name: "waf-ip-delete", ArgsUsage: "<dist-id> <ip-id>", Action: cdnDel("/cdn/v2/distributions/%s/waf/ipWhitelist/%s", 2, "distribution ID and IP ID required"), HideHelpCommand: true}

// Geo Restrictions
var cdnGeoRestrictionsCmd = cli.Command{Name: "geo-restrictions", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/geoRestrictions", 1, "distribution ID required"), HideHelpCommand: true}
var cdnGeoRestrictionsUpdateCmd = cli.Command{Name: "geo-restrictions-update", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/geoRestrictions", 1, "distribution ID required"), HideHelpCommand: true}

// Rate Limiting
var cdnRateLimitCmd = cli.Command{Name: "rate-limit", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/rateLimiting", 1, "distribution ID required"), HideHelpCommand: true}
var cdnRateLimitUpdateCmd = cli.Command{Name: "rate-limit-update", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/rateLimiting", 1, "distribution ID required"), HideHelpCommand: true}

// Metrics
func cdnMetricHandler(metricType string) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) < 1 {
			return fmt.Errorf("distribution ID required")
		}
		client, err := NewClient(cmd)
		if err != nil {
			return err
		}
		q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to"), "granularity": cmd.String("granularity")})
		res, err := client.Get(ctx, fmt.Sprintf("/cdn/v2/distributions/%s/metrics/%s%s", args[0], metricType, q))
		if err != nil {
			return err
		}
		return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
	}
}

var metricsFlags = []cli.Flag{&cli.StringFlag{Name: "from", Required: true}, &cli.StringFlag{Name: "to", Required: true}, &cli.StringFlag{Name: "granularity", Value: "1h"}}

var cdnMetricsBandwidthCmd = cli.Command{Name: "metrics-bandwidth", ArgsUsage: "<dist-id>", Flags: metricsFlags, Action: cdnMetricHandler("bandwidth"), HideHelpCommand: true}
var cdnMetricsRequestsCmd = cli.Command{Name: "metrics-requests", ArgsUsage: "<dist-id>", Flags: metricsFlags, Action: cdnMetricHandler("requests"), HideHelpCommand: true}
var cdnMetricsStatusCodesCmd = cli.Command{Name: "metrics-status-codes", ArgsUsage: "<dist-id>", Flags: metricsFlags, Action: cdnMetricHandler("statusCodes"), HideHelpCommand: true}
var cdnMetricsCacheHitCmd = cli.Command{Name: "metrics-cache-hit", ArgsUsage: "<dist-id>", Flags: metricsFlags, Action: cdnMetricHandler("cacheHitRatio"), HideHelpCommand: true}
var cdnMetricsEdgeLocationsCmd = cli.Command{Name: "metrics-edge-locations", ArgsUsage: "<dist-id>", Flags: metricsFlags, Action: cdnMetricHandler("edgeLocations"), HideHelpCommand: true}
var cdnMetricsOriginBandwidthCmd = cli.Command{Name: "metrics-origin-bandwidth", ArgsUsage: "<dist-id>", Flags: metricsFlags, Action: cdnMetricHandler("originBandwidth"), HideHelpCommand: true}

// Access Logs
var cdnAccessLogsCmd = cli.Command{Name: "access-logs", ArgsUsage: "<dist-id>", Flags: []cli.Flag{&cli.StringFlag{Name: "from"}, &cli.StringFlag{Name: "to"}}, Action: func(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()
	if len(args) < 1 {
		return fmt.Errorf("distribution ID required")
	}
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	q := BuildQueryString(map[string]string{"from": cmd.String("from"), "to": cmd.String("to")})
	res, err := client.Get(ctx, fmt.Sprintf("/cdn/v2/distributions/%s/accessLogs%s", args[0], q))
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

// Edge Locations
var cdnEdgeLocationsCmd = cli.Command{Name: "edge-locations", Usage: "List CDN edge locations", Action: func(ctx context.Context, cmd *cli.Command) error {
	client, err := NewClient(cmd)
	if err != nil {
		return err
	}
	res, err := client.Get(ctx, "/cdn/v2/edgeLocations")
	if err != nil {
		return err
	}
	return ShowResult(os.Stdout, res, cmd.Root().String("output"), cmd.Root().String("transform"))
}, HideHelpCommand: true}

// Token Auth
var cdnTokenAuthCmd = cli.Command{Name: "token-auth", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/tokenAuthentication", 1, "distribution ID required"), HideHelpCommand: true}
var cdnTokenAuthUpdateCmd = cli.Command{Name: "token-auth-update", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/tokenAuthentication", 1, "distribution ID required"), HideHelpCommand: true}

// Hotlink Protection
var cdnHotlinkProtectionCmd = cli.Command{Name: "hotlink-protection", ArgsUsage: "<dist-id>", Action: cdnCRUD("/cdn/v2/distributions/%s/hotlinkProtection", 1, "distribution ID required"), HideHelpCommand: true}
var cdnHotlinkProtectionUpdateCmd = cli.Command{Name: "hotlink-protection-update", ArgsUsage: "<dist-id>", Flags: payloadFlag, Action: cdnPut("/cdn/v2/distributions/%s/hotlinkProtection", 1, "distribution ID required"), HideHelpCommand: true}
