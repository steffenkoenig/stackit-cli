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
		Use:   fmt.Sprintf("delete %s", filterIdArg),
		Short: "Deletes a VPN BGP filter",
		Long:  "Deletes a VPN BGP filter. A filter which is still attached to a tunnel cannot be deleted.",
		Args:  args.SingleArg(filterIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete the BGP filter with the ID "yyy" of the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter delete yyy --gateway-id xxx",
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

			prompt := fmt.Sprintf("Are you sure you want to delete the BGP filter %q? (This cannot be undone)", filterLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err := req.Execute(); err != nil {
				return fmt.Errorf("delete vpn bgp filter: %w", err)
			}

			params.Printer.Outputf("Deleted BGP filter %q\n", filterLabel)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiDeleteGatewayBGPFilterRequest {
	return apiClient.DefaultAPI.DeleteGatewayBGPFilter(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId)
}
