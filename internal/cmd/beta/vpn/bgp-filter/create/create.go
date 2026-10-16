package create

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
	gatewayIdFlag   = "gateway-id"
	displayNameFlag = "display-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId   string
	DisplayName string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a VPN BGP filter",
		Long: `Creates a VPN BGP filter.

The filter is created without any rules. Since an empty filter denies every route, add rules with
'stackit beta vpn bgp-filter rule create' before attaching it to a tunnel.`,
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a VPN BGP filter named "my-filter" for the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter create --gateway-id xxx --display-name my-filter",
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

			gatewayLabel, err := vpnUtils.GetGatewayName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.GatewayId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get gateway name: %v", err)
			}
			if gatewayLabel == "" {
				gatewayLabel = model.GatewayId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a BGP filter for gateway %q?", gatewayLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create vpn bgp filter: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, gatewayLabel, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringValue(p, cmd, gatewayIdFlag),
		DisplayName:     flags.FlagToStringValue(p, cmd, displayNameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiCreateGatewayBGPFilterRequest {
	req := apiClient.DefaultAPI.CreateGatewayBGPFilter(ctx, model.ProjectId, model.Region, model.GatewayId)
	return req.CreateGatewayBGPFilterPayload(vpn.CreateGatewayBGPFilterPayload{
		DisplayName: model.DisplayName,
	})
}

func outputResult(p *print.Printer, outputFormat, gatewayLabel string, filter *vpn.BGPFilter) error {
	return p.OutputResult(outputFormat, filter, func() error {
		if filter == nil {
			p.Outputln("vpn bgp filter response is empty")
			return nil
		}
		p.Outputf("Created BGP filter %q for gateway %q.\nFilter ID: %s\n", filter.DisplayName, gatewayLabel, utils.PtrString(filter.Id))
		return nil
	})
}
