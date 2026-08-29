package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	vpnUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	connectionIdArg = "CONNECTION_ID"

	gatewayIdFlag = "gateway-id"

	displayNameFlag   = "display-name"
	enabledFlag       = "enabled"
	labelsFlag        = "labels"
	localSubnetsFlag  = "local-subnets"
	remoteSubnetsFlag = "remote-subnets"
	staticRoutesFlag  = "static-routes"

	tunnel1BgpInboundFilterIdFlag   = "tunnel1-bgp-inbound-filter-id"
	tunnel1BgpRemoteAsnFlag         = "tunnel1-bgp-remote-asn"
	tunnel1PeeringLocalAddressFlag  = "tunnel1-peering-local-address"
	tunnel1PeeringRemoteAddressFlag = "tunnel1-peering-remote-address"
	tunnel1Phase1RekeyTimeFlag      = "tunnel1-phase1-rekey-time"
	tunnel1Phase2RekeyTimeFlag      = "tunnel1-phase2-rekey-time"
	tunnel1PreSharedKeyFlag         = "tunnel1-pre-shared-key"
	tunnel1RemoteAddressFlag        = "tunnel1-remote-address"

	tunnel2BgpInboundFilterIdFlag   = "tunnel2-bgp-inbound-filter-id"
	tunnel2BgpRemoteAsnFlag         = "tunnel2-bgp-remote-asn"
	tunnel2PeeringLocalAddressFlag  = "tunnel2-peering-local-address"
	tunnel2PeeringRemoteAddressFlag = "tunnel2-peering-remote-address"
	tunnel2Phase1RekeyTimeFlag      = "tunnel2-phase1-rekey-time"
	tunnel2Phase2RekeyTimeFlag      = "tunnel2-phase2-rekey-time"
	tunnel2PreSharedKeyFlag         = "tunnel2-pre-shared-key"
	tunnel2RemoteAddressFlag        = "tunnel2-remote-address"
)

var (
	// tunnel 1
	tunnel1Phase1DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase1-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 1 Phase 1 DH Groups.\nThe Diffie-Hellman Group. Required, except if AEAD algorithms are selected.",
	)
	tunnel1Phase1EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase1-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Tunnel 1 Phase 1 Encryption Algorithms",
	)
	tunnel1Phase1IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase1-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Tunnel 1 Phase 1 Integrity Algorithms",
	)
	tunnel1Phase2DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase2-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 1 Phase 2 DH Groups",
	)
	tunnel1Phase2EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase2-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Tunnel 1 Phase 2 Encryption Algorithms",
	)
	tunnel1Phase2IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel1-phase2-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Tunnel 1 Phase 2 Integrity Algorithms",
	)
	tunnel1Phase2DpdActionFlag = flags.StringEnumFlag(
		"tunnel1-phase2-dpd-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfDpdActionEnumValues,
		"Tunnel 1 Phase 2 DPD Action.\nAction to perform for this CHILD_SA on DPD timeout.",
	)
	tunnel1Phase2StartActionFlag = flags.StringEnumFlag(
		"tunnel1-phase2-start-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfStartActionEnumValues,
		"Tunnel 1 Phase 2 Start Action.\nAction to perform after loading the connection configuration.",
	)
	// tunnel 2
	tunnel2Phase1DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase1-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 2 Phase 1 DH Groups.\nThe Diffie-Hellman Group. Required, except if AEAD algorithms are selected.",
	)
	tunnel2Phase1EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase1-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Tunnel 2 Phase 1 Encryption Algorithms",
	)
	tunnel2Phase1IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase1-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Tunnel 2 Phase 1 Integrity Algorithms",
	)
	tunnel2Phase2DhGroupsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase2-dh-groups",
		vpn.AllowedPhaseDhGroupsInnerEnumValues,
		"Tunnel 2 Phase 2 DH Groups",
	)
	tunnel2Phase2EncryptionAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase2-encryption-algorithms",
		vpn.AllowedPhaseEncryptionAlgorithmsInnerEnumValues,
		"Tunnel 2 Phase 2 Encryption Algorithms",
	)
	tunnel2Phase2IntegrityAlgorithmsFlag = flags.StringEnumSliceFlag(
		"tunnel2-phase2-integrity-algorithms",
		vpn.AllowedPhaseIntegrityAlgorithmsInnerEnumValues,
		"Tunnel 2 Phase 2 Integrity Algorithms",
	)
	tunnel2Phase2DpdActionFlag = flags.StringEnumFlag(
		"tunnel2-phase2-dpd-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfDpdActionEnumValues,
		"Tunnel 2 Phase 2 DPD Action.\nAction to perform for this CHILD_SA on DPD timeout.",
	)
	tunnel2Phase2StartActionFlag = flags.StringEnumFlag(
		"tunnel2-phase2-start-action",
		vpn.AllowedTunnelConfigurationPhase2AllOfStartActionEnumValues,
		"Tunnel 2 Phase 2 Start Action.\nAction to perform after loading the connection configuration.",
	)
)

type tunnelInputModel struct {
	BgpInboundFilterId         *string
	BgpRemoteAsn               *int64
	PeeringLocalAddress        *string
	PeeringRemoteAddress       *string
	Phase1DhGroups             []vpn.PhaseDhGroupsInner
	Phase1EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Phase1IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Phase1RekeyTime            *int32
	Phase2DhGroups             []vpn.PhaseDhGroupsInner
	Phase2EncryptionAlgorithms []vpn.PhaseEncryptionAlgorithmsInner
	Phase2IntegrityAlgorithms  []vpn.PhaseIntegrityAlgorithmsInner
	Phase2RekeyTime            *int32
	Phase2DpdAction            *vpn.TunnelConfigurationPhase2AllOfDpdAction
	Phase2StartAction          *vpn.TunnelConfigurationPhase2AllOfStartAction
	PreSharedKey               *string
	RemoteAddress              *string
}

func (t *tunnelInputModel) isEmpty() bool {
	return t.BgpInboundFilterId == nil && t.BgpRemoteAsn == nil &&
		t.PeeringLocalAddress == nil && t.PeeringRemoteAddress == nil &&
		t.Phase1DhGroups == nil && t.Phase1EncryptionAlgorithms == nil &&
		t.Phase1IntegrityAlgorithms == nil && t.Phase1RekeyTime == nil &&
		t.Phase2DhGroups == nil && t.Phase2EncryptionAlgorithms == nil &&
		t.Phase2IntegrityAlgorithms == nil && t.Phase2RekeyTime == nil &&
		t.Phase2DpdAction == nil && t.Phase2StartAction == nil &&
		t.PreSharedKey == nil && t.RemoteAddress == nil
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId    string
	ConnectionId string

	DisplayName   *string
	Enabled       *bool
	Labels        *map[string]string
	LocalSubnets  []string
	RemoteSubnets []string
	StaticRoutes  []string

	Tunnel1 tunnelInputModel
	Tunnel2 tunnelInputModel
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", connectionIdArg),
		Short: "Updates a VPN connection",
		Long: `Updates a VPN connection.

Only the attributes given as flags are changed; every other attribute keeps its current value.
Pre-shared keys are never returned by the API, so they are only changed when passed explicitly.`,
		Args: args.SingleArg(connectionIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Rename the VPN connection with the ID "zzz"`,
				"$ stackit beta vpn connection update zzz --gateway-id xxx --display-name my-connection"),
			examples.NewExample(
				`Attach the BGP filter "yyy" to the BGP session of tunnel 1`,
				"$ stackit beta vpn connection update zzz --gateway-id xxx --tunnel1-bgp-inbound-filter-id yyy"),
			examples.NewExample(
				`Disable the VPN connection with the ID "zzz"`,
				"$ stackit beta vpn connection update zzz --gateway-id xxx --enabled=false"),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, inputArgs)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// The update replaces the connection as a whole, so the current state is read first and
			// only the attributes given as flags are overlaid onto it.
			current, err := apiClient.DefaultAPI.GetGatewayConnection(ctx, model.ProjectId, model.Region, model.GatewayId, model.ConnectionId).Execute()
			if err != nil {
				return fmt.Errorf("read VPN connection: %w", err)
			}
			if current == nil {
				return fmt.Errorf("read VPN connection: response is empty")
			}

			connectionLabel := current.DisplayName
			if connectionLabel == "" {
				connectionLabel = model.ConnectionId
			}

			if err := validateBgpInboundFilters(model, current); err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to update the VPN connection %q?", connectionLabel)
			if err := p.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient, current)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update VPN connection: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, connectionLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Required: Gateway ID")
	cmd.Flags().String(displayNameFlag, "", "A user friendly name for the connection.")
	cmd.Flags().Bool(enabledFlag, true, "Enable the connection")
	cmd.Flags().StringToString(labelsFlag, nil, "Map of custom labels. Key and values must be a string with max 63 chars, start/end with alphanumeric. (example: foo=bar)")
	cmd.Flags().StringSlice(localSubnetsFlag, nil, "Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.")
	cmd.Flags().StringSlice(remoteSubnetsFlag, nil, "Defaults to 0.0.0.0/0 for Route-based VPN configurations. Mandatory for Policy-based.")
	cmd.Flags().StringSlice(staticRoutesFlag, nil, "Use this for route-based VPN.")

	cmd.Flags().Var(flags.UUIDFlag(), tunnel1BgpInboundFilterIdFlag, "Tunnel 1 BGP Inbound Filter ID.\nID of the BGP filter to apply to incoming routes from this tunnel's BGP neighbor. The filter must exist in the same gateway and the tunnel must have a BGP remote ASN.")
	cmd.Flags().Int64(tunnel1BgpRemoteAsnFlag, 0, "Tunnel 1 BGP Remote ASN.\nASN for private use (reserved by IANA), both 16Bit and 32Bit ranges are valid (RFC 6996).")
	cmd.Flags().String(tunnel1PeeringLocalAddressFlag, "", "Tunnel 1 Peering Local Address")
	cmd.Flags().String(tunnel1PeeringRemoteAddressFlag, "", "Tunnel 1 Peering Remote Address")
	tunnel1Phase1DhGroupsFlag.Register(cmd)
	tunnel1Phase1EncryptionAlgorithmsFlag.Register(cmd)
	tunnel1Phase1IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel1Phase1RekeyTimeFlag, 0, "Tunnel 1 Phase 1 Rekey Time.\nTime to schedule a IKE re-keying (in seconds).")
	tunnel1Phase2DhGroupsFlag.Register(cmd)
	tunnel1Phase2EncryptionAlgorithmsFlag.Register(cmd)
	tunnel1Phase2IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel1Phase2RekeyTimeFlag, 0, "Tunnel 1 Phase 2 Rekey Time.\nTime to schedule a Child SA re-keying (in seconds).")
	tunnel1Phase2DpdActionFlag.Register(cmd.Flags())
	tunnel1Phase2StartActionFlag.Register(cmd.Flags())
	cmd.Flags().String(tunnel1PreSharedKeyFlag, "", "Tunnel 1 Pre Shared Key.\nOnly changed when passed explicitly.")
	cmd.Flags().String(tunnel1RemoteAddressFlag, "", "Tunnel 1 Remote Address")

	cmd.Flags().Var(flags.UUIDFlag(), tunnel2BgpInboundFilterIdFlag, "Tunnel 2 BGP Inbound Filter ID.\nID of the BGP filter to apply to incoming routes from this tunnel's BGP neighbor. The filter must exist in the same gateway and the tunnel must have a BGP remote ASN.")
	cmd.Flags().Int64(tunnel2BgpRemoteAsnFlag, 0, "Tunnel 2 BGP Remote ASN")
	cmd.Flags().String(tunnel2PeeringLocalAddressFlag, "", "Tunnel 2 Peering Local Address")
	cmd.Flags().String(tunnel2PeeringRemoteAddressFlag, "", "Tunnel 2 Peering Remote Address")
	tunnel2Phase1DhGroupsFlag.Register(cmd)
	tunnel2Phase1EncryptionAlgorithmsFlag.Register(cmd)
	tunnel2Phase1IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel2Phase1RekeyTimeFlag, 0, "Tunnel 2 Phase 1 Rekey Time.\nTime to schedule a IKE re-keying (in seconds).")
	tunnel2Phase2DhGroupsFlag.Register(cmd)
	tunnel2Phase2EncryptionAlgorithmsFlag.Register(cmd)
	tunnel2Phase2IntegrityAlgorithmsFlag.Register(cmd)
	cmd.Flags().Int64(tunnel2Phase2RekeyTimeFlag, 0, "Tunnel 2 Phase 2 Rekey Time.\nTime to schedule a Child SA re-keying (in seconds).")
	tunnel2Phase2DpdActionFlag.Register(cmd.Flags())
	tunnel2Phase2StartActionFlag.Register(cmd.Flags())
	cmd.Flags().String(tunnel2PreSharedKeyFlag, "", "Tunnel 2 Pre Shared Key.\nOnly changed when passed explicitly.")
	cmd.Flags().String(tunnel2RemoteAddressFlag, "", "Tunnel 2 Remote Address")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringValue(p, cmd, gatewayIdFlag),
		ConnectionId:    inputArgs[0],

		DisplayName:   flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Enabled:       flags.FlagToBoolPointer(p, cmd, enabledFlag),
		Labels:        flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		LocalSubnets:  flags.FlagToStringSliceValue(p, cmd, localSubnetsFlag),
		RemoteSubnets: flags.FlagToStringSliceValue(p, cmd, remoteSubnetsFlag),
		StaticRoutes:  flags.FlagToStringSliceValue(p, cmd, staticRoutesFlag),

		Tunnel1: tunnelInputModel{
			BgpInboundFilterId:         flags.FlagToStringPointer(p, cmd, tunnel1BgpInboundFilterIdFlag),
			BgpRemoteAsn:               flags.FlagToInt64Pointer(p, cmd, tunnel1BgpRemoteAsnFlag),
			PeeringLocalAddress:        flags.FlagToStringPointer(p, cmd, tunnel1PeeringLocalAddressFlag),
			PeeringRemoteAddress:       flags.FlagToStringPointer(p, cmd, tunnel1PeeringRemoteAddressFlag),
			Phase1DhGroups:             tunnel1Phase1DhGroupsFlag.Get(),
			Phase1EncryptionAlgorithms: tunnel1Phase1EncryptionAlgorithmsFlag.Get(),
			Phase1IntegrityAlgorithms:  tunnel1Phase1IntegrityAlgorithmsFlag.Get(),
			Phase1RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel1Phase1RekeyTimeFlag),
			Phase2DhGroups:             tunnel1Phase2DhGroupsFlag.Get(),
			Phase2EncryptionAlgorithms: tunnel1Phase2EncryptionAlgorithmsFlag.Get(),
			Phase2IntegrityAlgorithms:  tunnel1Phase2IntegrityAlgorithmsFlag.Get(),
			Phase2RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel1Phase2RekeyTimeFlag),
			Phase2DpdAction:            tunnel1Phase2DpdActionFlag.Ptr(),
			Phase2StartAction:          tunnel1Phase2StartActionFlag.Ptr(),
			PreSharedKey:               flags.FlagToStringPointer(p, cmd, tunnel1PreSharedKeyFlag),
			RemoteAddress:              flags.FlagToStringPointer(p, cmd, tunnel1RemoteAddressFlag),
		},

		Tunnel2: tunnelInputModel{
			BgpInboundFilterId:         flags.FlagToStringPointer(p, cmd, tunnel2BgpInboundFilterIdFlag),
			BgpRemoteAsn:               flags.FlagToInt64Pointer(p, cmd, tunnel2BgpRemoteAsnFlag),
			PeeringLocalAddress:        flags.FlagToStringPointer(p, cmd, tunnel2PeeringLocalAddressFlag),
			PeeringRemoteAddress:       flags.FlagToStringPointer(p, cmd, tunnel2PeeringRemoteAddressFlag),
			Phase1DhGroups:             tunnel2Phase1DhGroupsFlag.Get(),
			Phase1EncryptionAlgorithms: tunnel2Phase1EncryptionAlgorithmsFlag.Get(),
			Phase1IntegrityAlgorithms:  tunnel2Phase1IntegrityAlgorithmsFlag.Get(),
			Phase1RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel2Phase1RekeyTimeFlag),
			Phase2DhGroups:             tunnel2Phase2DhGroupsFlag.Get(),
			Phase2EncryptionAlgorithms: tunnel2Phase2EncryptionAlgorithmsFlag.Get(),
			Phase2IntegrityAlgorithms:  tunnel2Phase2IntegrityAlgorithmsFlag.Get(),
			Phase2RekeyTime:            flags.FlagToInt32Pointer(p, cmd, tunnel2Phase2RekeyTimeFlag),
			Phase2DpdAction:            tunnel2Phase2DpdActionFlag.Ptr(),
			Phase2StartAction:          tunnel2Phase2StartActionFlag.Ptr(),
			PreSharedKey:               flags.FlagToStringPointer(p, cmd, tunnel2PreSharedKeyFlag),
			RemoteAddress:              flags.FlagToStringPointer(p, cmd, tunnel2RemoteAddressFlag),
		},
	}

	if model.DisplayName == nil && model.Enabled == nil && model.Labels == nil &&
		model.LocalSubnets == nil && model.RemoteSubnets == nil && model.StaticRoutes == nil &&
		model.Tunnel1.isEmpty() && model.Tunnel2.isEmpty() {
		return nil, &cliErr.EmptyUpdateError{}
	}

	warnDeprecatedIntegrityAlgorithms(p, &model)

	p.DebugInputModel(model)
	return &model, nil
}

// validateBgpInboundFilters rejects attaching a BGP filter to a tunnel which ends up without a BGP
// session, since a filter only applies to routes learned over BGP. The remote ASN may come either
// from a flag or from the tunnel's current configuration.
func validateBgpInboundFilters(model *inputModel, current *vpn.ConnectionResponse) error {
	tunnels := []struct {
		model               *tunnelInputModel
		current             *vpn.TunnelConfiguration
		remoteAsnFlag       string
		inboundFilterIdFlag string
	}{
		{&model.Tunnel1, &current.Tunnel1, tunnel1BgpRemoteAsnFlag, tunnel1BgpInboundFilterIdFlag},
		{&model.Tunnel2, &current.Tunnel2, tunnel2BgpRemoteAsnFlag, tunnel2BgpInboundFilterIdFlag},
	}
	for _, tunnel := range tunnels {
		if tunnel.model.BgpInboundFilterId == nil {
			continue
		}
		if tunnel.model.BgpRemoteAsn != nil || tunnel.current.Bgp != nil {
			continue
		}
		return &cliErr.DependingFlagIsMissing{
			MissingFlag: tunnel.remoteAsnFlag,
			SetFlag:     tunnel.inboundFilterIdFlag,
		}
	}
	return nil
}

// warnDeprecatedIntegrityAlgorithms warns about every deprecated integrity algorithm configured for
// either tunnel, so users can migrate to a stronger algorithm.
func warnDeprecatedIntegrityAlgorithms(p *print.Printer, model *inputModel) {
	vpnUtils.WarnDeprecatedIntegrityAlgorithms(p, tunnel1Phase1IntegrityAlgorithmsFlag.Name(), model.Tunnel1.Phase1IntegrityAlgorithms)
	vpnUtils.WarnDeprecatedIntegrityAlgorithms(p, tunnel1Phase2IntegrityAlgorithmsFlag.Name(), model.Tunnel1.Phase2IntegrityAlgorithms)
	vpnUtils.WarnDeprecatedIntegrityAlgorithms(p, tunnel2Phase1IntegrityAlgorithmsFlag.Name(), model.Tunnel2.Phase1IntegrityAlgorithms)
	vpnUtils.WarnDeprecatedIntegrityAlgorithms(p, tunnel2Phase2IntegrityAlgorithmsFlag.Name(), model.Tunnel2.Phase2IntegrityAlgorithms)
}

// mergeTunnelConfiguration overlays the attributes given as flags onto the tunnel's current
// configuration. Attributes which are not given keep the value they currently have.
func mergeTunnelConfiguration(model *tunnelInputModel, current *vpn.TunnelConfiguration) vpn.TunnelConfiguration {
	// Copied so the caller's configuration is left untouched.
	tunnel := *current

	if model.RemoteAddress != nil {
		tunnel.RemoteAddress = *model.RemoteAddress
	}
	// The API never returns the pre-shared key, so it is only sent when passed explicitly.
	if model.PreSharedKey != nil {
		tunnel.PreSharedKey = model.PreSharedKey
	} else {
		tunnel.PreSharedKey = nil
	}

	if model.BgpRemoteAsn != nil || model.BgpInboundFilterId != nil {
		bgp := vpn.BGPTunnelConfig{}
		if tunnel.Bgp != nil {
			bgp = *tunnel.Bgp
		}
		if model.BgpRemoteAsn != nil {
			bgp.RemoteAsn = *model.BgpRemoteAsn
		}
		if model.BgpInboundFilterId != nil {
			bgp.InboundFilterId = *vpn.NewNullableString(model.BgpInboundFilterId)
		}
		tunnel.Bgp = &bgp
	}

	if model.PeeringLocalAddress != nil || model.PeeringRemoteAddress != nil {
		peering := vpn.PeeringConfig{}
		if tunnel.Peering != nil {
			peering = *tunnel.Peering
		}
		if model.PeeringLocalAddress != nil {
			peering.LocalAddress = model.PeeringLocalAddress
		}
		if model.PeeringRemoteAddress != nil {
			peering.RemoteAddress = model.PeeringRemoteAddress
		}
		tunnel.Peering = &peering
	}

	if model.Phase1DhGroups != nil {
		tunnel.Phase1.DhGroups = model.Phase1DhGroups
	}
	if model.Phase1EncryptionAlgorithms != nil {
		tunnel.Phase1.EncryptionAlgorithms = model.Phase1EncryptionAlgorithms
	}
	if model.Phase1IntegrityAlgorithms != nil {
		tunnel.Phase1.IntegrityAlgorithms = model.Phase1IntegrityAlgorithms
	}
	if model.Phase1RekeyTime != nil {
		tunnel.Phase1.RekeyTime = model.Phase1RekeyTime
	}

	if model.Phase2DhGroups != nil {
		tunnel.Phase2.DhGroups = model.Phase2DhGroups
	}
	if model.Phase2EncryptionAlgorithms != nil {
		tunnel.Phase2.EncryptionAlgorithms = model.Phase2EncryptionAlgorithms
	}
	if model.Phase2IntegrityAlgorithms != nil {
		tunnel.Phase2.IntegrityAlgorithms = model.Phase2IntegrityAlgorithms
	}
	if model.Phase2RekeyTime != nil {
		tunnel.Phase2.RekeyTime = model.Phase2RekeyTime
	}
	if model.Phase2DpdAction != nil {
		tunnel.Phase2.DpdAction = model.Phase2DpdAction
	}
	if model.Phase2StartAction != nil {
		tunnel.Phase2.StartAction = model.Phase2StartAction
	}

	return tunnel
}

// buildRequest overlays the attributes given as flags onto the current connection. Attributes which
// are not given keep the value they currently have, since the API replaces the connection as a whole.
func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient, current *vpn.ConnectionResponse) vpn.ApiUpdateGatewayConnectionRequest {
	payload := vpn.UpdateGatewayConnectionPayload{
		DisplayName:   current.DisplayName,
		Enabled:       current.Enabled,
		Labels:        current.Labels,
		LocalSubnets:  current.LocalSubnets,
		RemoteSubnets: current.RemoteSubnets,
		StaticRoutes:  current.StaticRoutes,
	}

	if model.DisplayName != nil {
		payload.DisplayName = *model.DisplayName
	}
	if model.Enabled != nil {
		payload.Enabled = model.Enabled
	}
	if model.Labels != nil {
		payload.Labels = model.Labels
	}
	if model.LocalSubnets != nil {
		payload.LocalSubnets = model.LocalSubnets
	}
	if model.RemoteSubnets != nil {
		payload.RemoteSubnets = model.RemoteSubnets
	}
	if model.StaticRoutes != nil {
		payload.StaticRoutes = model.StaticRoutes
	}

	payload.Tunnel1 = mergeTunnelConfiguration(&model.Tunnel1, &current.Tunnel1)
	payload.Tunnel2 = mergeTunnelConfiguration(&model.Tunnel2, &current.Tunnel2)

	req := apiClient.DefaultAPI.UpdateGatewayConnection(ctx, model.ProjectId, model.Region, model.GatewayId, model.ConnectionId)
	return req.UpdateGatewayConnectionPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, connectionLabel string, resp *vpn.ConnectionResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			p.Outputln("VPN connection response is empty")
			return nil
		}
		p.Outputf("Updated VPN connection %q.\n", connectionLabel)
		return nil
	})
}
