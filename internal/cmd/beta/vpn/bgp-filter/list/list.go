package list

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	gatewayIdFlag = "gateway-id"
	limitFlag     = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
	Limit     *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all VPN BGP filters of a gateway",
		Long:  "Lists all VPN BGP filters of a gateway.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all BGP filters of the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter list --gateway-id xxx",
			),
			examples.NewExample(
				`List up to 4 BGP filters of the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter list --gateway-id xxx --limit 4",
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
				return fmt.Errorf("list vpn bgp filters: %w", err)
			}

			filters := resp.BgpFilters
			// Truncate output
			if model.Limit != nil && len(filters) > int(*model.Limit) {
				filters = filters[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, model.GatewayId, filters)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag)
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
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiListGatewayBGPFiltersRequest {
	return apiClient.DefaultAPI.ListGatewayBGPFilters(ctx, model.ProjectId, model.Region, model.GatewayId)
}

func outputResult(p *print.Printer, outputFormat, gatewayId string, filters []vpn.BGPFilter) error {
	return p.OutputResult(outputFormat, filters, func() error {
		if len(filters) == 0 {
			p.Info("No BGP filters found for gateway %q\n", gatewayId)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME")
		for _, filter := range filters {
			table.AddRow(utils.PtrString(filter.Id), filter.DisplayName)
		}
		p.Outputln(table.Render())
		return nil
	})
}
