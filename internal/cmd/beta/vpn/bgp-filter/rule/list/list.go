package list

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	gatewayIdFlag = "gateway-id"
	filterIdFlag  = "filter-id"
	limitFlag     = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
	FilterId  string
	Limit     *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all rules of a VPN BGP filter",
		Long:  "Lists all rules of a VPN BGP filter, ordered by their sequence.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all rules of the BGP filter with the ID "yyy" of the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter rule list --gateway-id xxx --filter-id yyy",
			),
			examples.NewExample(
				`List up to 4 rules of the BGP filter with the ID "yyy"`,
				"$ stackit beta vpn bgp-filter rule list --gateway-id xxx --filter-id yyy --limit 4",
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list vpn bgp filter rules: %w", err)
			}

			rules := resp.Rules
			// Truncate output
			if model.Limit != nil && len(rules) > int(*model.Limit) {
				rules = rules[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, model.FilterId, rules)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().Var(flags.UUIDFlag(), filterIdFlag, "BGP filter ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag, filterIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringValue(p, cmd, gatewayIdFlag),
		FilterId:        flags.FlagToStringValue(p, cmd, filterIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiListGatewayBGPFilterRulesRequest {
	return apiClient.DefaultAPI.ListGatewayBGPFilterRules(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId)
}

func outputResult(p *print.Printer, outputFormat, filterId string, rules []vpn.BGPFilterRule) error {
	return p.OutputResult(outputFormat, rules, func() error {
		if len(rules) == 0 {
			p.Info("No rules found for BGP filter %q\n", filterId)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("SEQUENCE", "ID", "ACTION", "MATCH", "LOCAL PREFERENCE")
		for _, rule := range rules {
			localPreference := ""
			if rule.Set != nil && rule.Set.LocalPreference != nil {
				localPreference = utils.PtrString(rule.Set.LocalPreference)
			}
			table.AddRow(
				utils.PtrString(rule.Sequence),
				utils.PtrString(rule.Id),
				rule.Action,
				summarizeMatch(rule.Match),
				localPreference,
			)
		}
		p.Outputln(table.Render())
		return nil
	})
}

// summarizeMatch renders the match criteria of a rule as a single table cell.
func summarizeMatch(match *vpn.BGPFilterRuleMatch) string {
	if match == nil {
		return "any"
	}

	var criteria []string
	if len(match.Prefixes) > 0 {
		criteria = append(criteria, fmt.Sprintf("prefixes: %s", strings.Join(match.Prefixes, ", ")))
	}
	if match.MinPrefixLength != nil {
		criteria = append(criteria, fmt.Sprintf("min prefix length: %s", utils.PtrString(match.MinPrefixLength)))
	}
	if match.MaxPrefixLength != nil {
		criteria = append(criteria, fmt.Sprintf("max prefix length: %s", utils.PtrString(match.MaxPrefixLength)))
	}
	if match.Peer != nil {
		criteria = append(criteria, fmt.Sprintf("peer: %s", utils.PtrString(match.Peer)))
	}
	if match.FirstASN != nil {
		criteria = append(criteria, fmt.Sprintf("first ASN: %s", utils.PtrString(match.FirstASN)))
	}
	if len(match.AsPathContainsAny) > 0 {
		asns := make([]string, 0, len(match.AsPathContainsAny))
		for _, asn := range match.AsPathContainsAny {
			asns = append(asns, strconv.FormatInt(asn, 10))
		}
		criteria = append(criteria, fmt.Sprintf("AS-PATH contains any: %s", strings.Join(asns, ", ")))
	}
	if len(match.Communities) > 0 {
		criteria = append(criteria, fmt.Sprintf("communities: %s", strings.Join(match.Communities, ", ")))
	}

	if len(criteria) == 0 {
		return "any"
	}
	return strings.Join(criteria, "\n")
}
