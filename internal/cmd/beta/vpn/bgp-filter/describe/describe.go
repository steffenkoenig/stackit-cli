package describe

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
	filterIdArg = "FILTER_ID"

	gatewayIdFlag = "gateway-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
	FilterId  string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", filterIdArg),
		Short: "Shows details of a VPN BGP filter",
		Long:  "Shows details of a VPN BGP filter.",
		Args:  args.SingleArg(filterIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of the BGP filter with the ID "yyy" of the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter describe yyy --gateway-id xxx",
			),
			examples.NewExample(
				`Show details of the BGP filter with the ID "yyy" of the gateway with the ID "xxx" in JSON format`,
				"$ stackit beta vpn bgp-filter describe yyy --gateway-id xxx --output-format json",
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
				return fmt.Errorf("describe vpn bgp filter: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.FilterId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag)
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
		FilterId:        inputArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiGetGatewayBGPFilterRequest {
	return apiClient.DefaultAPI.GetGatewayBGPFilter(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId)
}

func outputResult(p *print.Printer, outputFormat, filterId string, filter *vpn.BGPFilter) error {
	return p.OutputResult(outputFormat, filter, func() error {
		if filter == nil {
			p.Outputf("bgp filter %q not found\n", filterId)
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("BGP Filter")
		table.AddRow("ID", utils.PtrString(filter.Id))
		table.AddSeparator()
		table.AddRow("NAME", filter.DisplayName)

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
