package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	vpnUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/utils"
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
		Use:   fmt.Sprintf("delete %s", ruleIdArg),
		Short: "Deletes a VPN BGP filter rule",
		Long: `Deletes a VPN BGP filter rule.

Since an implicit DENY follows the last rule of a filter, removing the last PERMIT rule of a filter
denies every route.`,
		Args: args.SingleArg(ruleIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete the rule with the ID "zzz" of the BGP filter with the ID "yyy"`,
				"$ stackit beta vpn bgp-filter rule delete zzz --gateway-id xxx --filter-id yyy",
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

			filterLabel, err := vpnUtils.GetBGPFilterName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.GatewayId, model.FilterId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get bgp filter name: %v", err)
			}
			if filterLabel == "" {
				filterLabel = model.FilterId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the rule %q of the BGP filter %q? (This cannot be undone)", model.RuleId, filterLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err := req.Execute(); err != nil {
				return fmt.Errorf("delete vpn bgp filter rule: %w", err)
			}

			params.Printer.Outputf("Deleted rule %q of BGP filter %q\n", model.RuleId, filterLabel)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiDeleteGatewayBGPFilterRuleRequest {
	return apiClient.DefaultAPI.DeleteGatewayBGPFilterRule(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId, model.RuleId)
}
