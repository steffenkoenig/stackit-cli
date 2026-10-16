package describe

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
	ruleIdArg = "RULE_ID"

	gatewayIdFlag = "gateway-id"
	filterIdFlag  = "filter-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
	FilterId  string
	RuleId    string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", ruleIdArg),
		Short: "Shows details of a VPN BGP filter rule",
		Long:  "Shows details of a VPN BGP filter rule.",
		Args:  args.SingleArg(ruleIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of the rule with the ID "zzz" of the BGP filter with the ID "yyy"`,
				"$ stackit beta vpn bgp-filter rule describe zzz --gateway-id xxx --filter-id yyy",
			),
			examples.NewExample(
				`Show details of the rule with the ID "zzz" of the BGP filter with the ID "yyy" in JSON format`,
				"$ stackit beta vpn bgp-filter rule describe zzz --gateway-id xxx --filter-id yyy --output-format json",
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
				return fmt.Errorf("describe vpn bgp filter rule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.RuleId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().Var(flags.UUIDFlag(), filterIdFlag, "BGP filter ID")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag, filterIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringValue(p, cmd, gatewayIdFlag),
		FilterId:        flags.FlagToStringValue(p, cmd, filterIdFlag),
		RuleId:          inputArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiGetGatewayBGPFilterRuleRequest {
	return apiClient.DefaultAPI.GetGatewayBGPFilterRule(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId, model.RuleId)
}

func outputResult(p *print.Printer, outputFormat, ruleId string, rule *vpn.BGPFilterRule) error {
	return p.OutputResult(outputFormat, rule, func() error {
		if rule == nil {
			p.Outputf("bgp filter rule %q not found\n", ruleId)
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("BGP Filter Rule")
		table.AddRow("ID", utils.PtrString(rule.Id))
		table.AddSeparator()
		table.AddRow("SEQUENCE", utils.PtrString(rule.Sequence))
		table.AddSeparator()
		table.AddRow("ACTION", rule.Action)

		if match := rule.Match; match != nil {
			if len(match.Prefixes) > 0 {
				table.AddSeparator()
				table.AddRow("Match Prefixes", strings.Join(match.Prefixes, "\n"))
			}
			if match.MinPrefixLength != nil {
				table.AddSeparator()
				table.AddRow("Match Min Prefix Length", utils.PtrString(match.MinPrefixLength))
			}
			if match.MaxPrefixLength != nil {
				table.AddSeparator()
				table.AddRow("Match Max Prefix Length", utils.PtrString(match.MaxPrefixLength))
			}
			if match.Peer != nil {
				table.AddSeparator()
				table.AddRow("Match Peer", utils.PtrString(match.Peer))
			}
			if match.FirstASN != nil {
				table.AddSeparator()
				table.AddRow("Match First ASN", utils.PtrString(match.FirstASN))
			}
			if len(match.AsPathContainsAny) > 0 {
				table.AddSeparator()
				table.AddRow("Match AS-PATH Contains Any", joinInt64s(match.AsPathContainsAny, "\n"))
			}
			if len(match.Communities) > 0 {
				table.AddSeparator()
				table.AddRow("Match Communities", strings.Join(match.Communities, "\n"))
			}
		}

		if rule.Set != nil && rule.Set.LocalPreference != nil {
			table.AddSeparator()
			table.AddRow("Set Local Preference", utils.PtrString(rule.Set.LocalPreference))
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}

func joinInt64s(values []int64, sep string) string {
	formatted := make([]string, 0, len(values))
	for _, value := range values {
		formatted = append(formatted, strconv.FormatInt(value, 10))
	}
	return strings.Join(formatted, sep)
}
