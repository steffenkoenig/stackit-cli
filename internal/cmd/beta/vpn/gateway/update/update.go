package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	gatewayIdArg = "GATEWAY_ID"

	availabilityZoneTunnel1Flag     = "availability-zone-tunnel-1"
	availabilityZoneTunnel2Flag     = "availability-zone-tunnel-2"
	bgpLocalAsnFlag                 = "bgp-local-asn"
	bgpOverrideAdvertisedRoutesFlag = "bgp-override-advertised-routes"
	nameFlag                        = "name"
	labelsFlag                      = "labels"
	networkRoutingTableIdFlag       = "network-routing-table-id"
	planIdFlag                      = "plan-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string

	AvailabilityZoneTunnel1     *string
	AvailabilityZoneTunnel2     *string
	BgpLocalAsn                 *int64
	BgpOverrideAdvertisedRoutes []string
	Name                        *string
	Labels                      *map[string]string
	NetworkRoutingTableId       *string
	PlanId                      *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", gatewayIdArg),
		Short: "Updates a vpn gateway",
		Long: `Updates a vpn gateway.

Only the attributes given as flags are changed; every other attribute keeps its current value.
The routing type of a gateway can only be set at creation and is therefore not updatable. The same
applies to the predefined network prefix.`,
		Args: args.SingleArg(gatewayIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Rename the vpn gateway with the ID "xxx" to "yyy"`,
				"$ stackit beta vpn gateway update xxx --name yyy",
			),
			examples.NewExample(
				`Update the labels of the vpn gateway with the ID "xxx"`,
				"$ stackit beta vpn gateway update xxx --labels foo=bar,x=y",
			),
			examples.NewExample(
				`Change the advertised BGP routes of the vpn gateway with the ID "xxx"`,
				"$ stackit beta vpn gateway update xxx --bgp-local-asn 64512 --bgp-override-advertised-routes 10.0.0.0/24,10.1.0.0/24",
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

			// The update replaces the gateway as a whole, so the current state is read first and
			// only the attributes given as flags are overlaid onto it.
			current, err := apiClient.DefaultAPI.GetGateway(ctx, model.ProjectId, model.Region, model.GatewayId).Execute()
			if err != nil {
				return fmt.Errorf("read vpn gateway: %w", err)
			}
			if current == nil {
				return fmt.Errorf("read vpn gateway: response is empty")
			}

			gatewayLabel := current.DisplayName
			if gatewayLabel == "" {
				gatewayLabel = model.GatewayId
			}

			prompt := fmt.Sprintf("Are you sure you want to update the vpn gateway %q?", gatewayLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient, current)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpn gateway: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating vpn gateway", func() error {
					_, err = wait.UpdateGatewayWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.GatewayId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("waiting for vpn gateway update: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, gatewayLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(availabilityZoneTunnel1Flag, "", "Availability Zone of Tunnel 1")
	cmd.Flags().String(availabilityZoneTunnel2Flag, "", "Availability Zone of Tunnel 2")
	cmd.Flags().Int64(bgpLocalAsnFlag, 0, "ASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996)")
	cmd.Flags().StringArray(bgpOverrideAdvertisedRoutesFlag, nil, "A list of IPv4 Prefixes to advertise via BGP")
	cmd.Flags().String(nameFlag, "", "Gateway name")
	cmd.Flags().StringToString(labelsFlag, nil, "Labels in key=value format, separated by commas")
	cmd.Flags().Var(flags.UUIDFlag(), networkRoutingTableIdFlag, "Custom routing table ID for the vpn gateway")
	cmd.Flags().String(planIdFlag, "", "Plan ID")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	bgpLocalAsn := flags.FlagToInt64Pointer(p, cmd, bgpLocalAsnFlag)
	if bgpLocalAsn != nil && *bgpLocalAsn < 0 {
		return nil, &errors.FlagValidationError{
			Flag:    bgpLocalAsnFlag,
			Details: "must be a positive integer",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       inputArgs[0],

		AvailabilityZoneTunnel1:     flags.FlagToStringPointer(p, cmd, availabilityZoneTunnel1Flag),
		AvailabilityZoneTunnel2:     flags.FlagToStringPointer(p, cmd, availabilityZoneTunnel2Flag),
		BgpLocalAsn:                 bgpLocalAsn,
		BgpOverrideAdvertisedRoutes: flags.FlagToStringArrayValue(p, cmd, bgpOverrideAdvertisedRoutesFlag),
		Name:                        flags.FlagToStringPointer(p, cmd, nameFlag),
		Labels:                      flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		NetworkRoutingTableId:       flags.FlagToStringPointer(p, cmd, networkRoutingTableIdFlag),
		PlanId:                      flags.FlagToStringPointer(p, cmd, planIdFlag),
	}

	if isEmpty(&model) {
		return nil, &errors.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return &model, nil
}

// isEmpty reports whether the update would not change anything.
func isEmpty(model *inputModel) bool {
	return model.AvailabilityZoneTunnel1 == nil &&
		model.AvailabilityZoneTunnel2 == nil &&
		model.BgpLocalAsn == nil &&
		model.BgpOverrideAdvertisedRoutes == nil &&
		model.Name == nil &&
		model.Labels == nil &&
		model.NetworkRoutingTableId == nil &&
		model.PlanId == nil
}

// buildRequest overlays the attributes given as flags onto the current gateway. Attributes which
// are not given keep the value they currently have, since the API replaces the gateway as a whole.
func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient, current *vpn.GatewayResponse) vpn.ApiUpdateGatewayRequest {
	payload := vpn.UpdateGatewayPayload{
		AvailabilityZones: vpn.UpdateGatewayPayloadAvailabilityZones{
			Tunnel1: current.AvailabilityZones.Tunnel1,
			Tunnel2: current.AvailabilityZones.Tunnel2,
		},
		Bgp:           current.Bgp,
		DisplayName:   current.DisplayName,
		Labels:        current.Labels,
		NetworkConfig: current.NetworkConfig,
		PlanId:        current.PlanId,
		// The routing type can only be set at creation, so the current value is always sent back.
		RoutingType: current.RoutingType,
	}

	if model.AvailabilityZoneTunnel1 != nil {
		payload.AvailabilityZones.Tunnel1 = *model.AvailabilityZoneTunnel1
	}
	if model.AvailabilityZoneTunnel2 != nil {
		payload.AvailabilityZones.Tunnel2 = *model.AvailabilityZoneTunnel2
	}
	if model.Name != nil {
		payload.DisplayName = *model.Name
	}
	if model.Labels != nil {
		payload.Labels = model.Labels
	}
	if model.PlanId != nil {
		payload.PlanId = *model.PlanId
	}

	if model.BgpLocalAsn != nil || model.BgpOverrideAdvertisedRoutes != nil {
		bgp := vpn.BGPGatewayConfig{}
		if payload.Bgp != nil {
			bgp = *payload.Bgp
		}
		if model.BgpLocalAsn != nil {
			bgp.LocalAsn = *model.BgpLocalAsn
		}
		if model.BgpOverrideAdvertisedRoutes != nil {
			bgp.OverrideAdvertisedRoutes = model.BgpOverrideAdvertisedRoutes
		}
		payload.Bgp = &bgp
	}

	if model.NetworkRoutingTableId != nil {
		networkConfig := vpn.NetworkConfig{}
		if payload.NetworkConfig != nil {
			networkConfig = *payload.NetworkConfig
		}
		networkConfig.RoutingTableId = model.NetworkRoutingTableId
		payload.NetworkConfig = &networkConfig
	}

	req := apiClient.DefaultAPI.UpdateGateway(ctx, model.ProjectId, model.Region, model.GatewayId)
	return req.UpdateGatewayPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, async bool, gatewayLabel string, gateway *vpn.GatewayResponse) error {
	return p.OutputResult(outputFormat, gateway, func() error {
		if gateway == nil {
			p.Outputln("vpn gateway response is empty")
			return nil
		}
		operation := "Updated"
		if async {
			operation = "Triggered update of"
		}
		p.Outputf("%s vpn gateway %q.\n", operation, gatewayLabel)
		return nil
	})
}
