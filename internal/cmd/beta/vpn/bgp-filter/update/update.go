package update

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
	filterIdArg = "FILTER_ID"

	gatewayIdFlag   = "gateway-id"
	displayNameFlag = "display-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId   string
	FilterId    string
	DisplayName string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", filterIdArg),
		Short: "Updates a VPN BGP filter",
		Long:  "Updates the metadata of a VPN BGP filter. The rules of the filter are managed via 'stackit beta vpn bgp-filter rule'.",
		Args:  args.SingleArg(filterIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Rename the BGP filter with the ID "yyy" of the gateway with the ID "xxx" to "my-filter"`,
				"$ stackit beta vpn bgp-filter update yyy --gateway-id xxx --display-name my-filter",
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

			prompt := fmt.Sprintf("Are you sure you want to update the BGP filter %q?", filterLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpn bgp filter: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, filterLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().String(displayNameFlag, "", "A user friendly name for the filter")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag, displayNameFlag)
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
		DisplayName:     flags.FlagToStringValue(p, cmd, displayNameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiUpdateGatewayBGPFilterRequest {
	req := apiClient.DefaultAPI.UpdateGatewayBGPFilter(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId)
	return req.UpdateGatewayBGPFilterPayload(vpn.UpdateGatewayBGPFilterPayload{
		DisplayName: model.DisplayName,
	})
}

func outputResult(p *print.Printer, outputFormat, filterLabel string, filter *vpn.BGPFilter) error {
	return p.OutputResult(outputFormat, filter, func() error {
		if filter == nil {
			p.Outputln("vpn bgp filter response is empty")
			return nil
		}
		p.Outputf("Updated BGP filter %q.\n", filterLabel)
		return nil
	})
}
