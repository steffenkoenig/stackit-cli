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
	gatewayIdFlag = "gateway-id"
	filterIdFlag  = "filter-id"

	sequenceFlag = "sequence"

	matchAsPathContainsAnyFlag = "match-as-path-contains-any"
	matchCommunitiesFlag       = "match-communities"
	matchFirstAsnFlag          = "match-first-asn"
	matchMaxPrefixLengthFlag   = "match-max-prefix-length"
	matchMinPrefixLengthFlag   = "match-min-prefix-length"
	matchPeerFlag              = "match-peer"
	matchPrefixesFlag          = "match-prefixes"

	setLocalPreferenceFlag = "set-local-preference"
)

var actionFlag = flags.StringEnumFlag(
	"action",
	vpn.AllowedCreateGatewayBGPFilterRulePayloadActionEnumValues,
	"The action to take if a route matches all criteria of the rule",
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
	FilterId  string

	Action   vpn.CreateGatewayBGPFilterRulePayloadAction
	Sequence *int32

	MatchAsPathContainsAny []int64
	MatchCommunities       []string
	MatchFirstAsn          *int64
	MatchMaxPrefixLength   *int32
	MatchMinPrefixLength   *int32
	MatchPeer              *string
	MatchPrefixes          []string

	SetLocalPreference *int32
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a VPN BGP filter rule",
		Long: `Creates a VPN BGP filter rule.

All match criteria which are set must match for the rule to apply (logical AND). If no match
criterion is set, the rule matches every route. When the sequence is omitted, the rule is appended
after the last existing rule of the filter.`,
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Permit every route received by the BGP filter with the ID "yyy" of the gateway with the ID "xxx"`,
				"$ stackit beta vpn bgp-filter rule create --gateway-id xxx --filter-id yyy --action PERMIT",
			),
			examples.NewExample(
				`Deny the routes 10.0.0.0/16 and 192.168.0.0/16 as the first rule of the filter`,
				"$ stackit beta vpn bgp-filter rule create --gateway-id xxx --filter-id yyy --action DENY --sequence 5 --match-prefixes 10.0.0.0/16,192.168.0.0/16",
			),
			examples.NewExample(
				`Permit the routes advertised by the peer 192.168.0.1 whose AS-PATH starts with the ASN 65000 and prefer them over other routes`,
				"$ stackit beta vpn bgp-filter rule create --gateway-id xxx --filter-id yyy --action PERMIT --match-peer 192.168.0.1 --match-first-asn 65000 --set-local-preference 200",
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

			prompt := fmt.Sprintf("Are you sure you want to create a rule for the BGP filter %q?", filterLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create vpn bgp filter rule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, filterLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().Var(flags.UUIDFlag(), filterIdFlag, "BGP filter ID")
	actionFlag.Register(cmd.Flags())
	cmd.Flags().Int32(sequenceFlag, 0, "The evaluation order of the rule. Lower numbers are evaluated first and the number must be unique within the filter. Defaults to the highest existing sequence plus 10.")

	cmd.Flags().Int64Slice(matchAsPathContainsAnyFlag, nil, "Match routes whose AS-PATH contains any of the given ASNs")
	cmd.Flags().StringSlice(matchCommunitiesFlag, nil, "Match routes carrying any of the given BGP standard communities in the 'asn:value' format (RFC 1997)")
	cmd.Flags().Int64(matchFirstAsnFlag, 0, "Match routes whose first ASN in the AS-PATH, i.e. the immediate neighbor, equals the given ASN")
	cmd.Flags().Int32(matchMaxPrefixLengthFlag, 0, "Match routes with a subnet mask length of at most this value")
	cmd.Flags().Int32(matchMinPrefixLengthFlag, 0, "Match routes with a subnet mask length of at least this value")
	cmd.Flags().String(matchPeerFlag, "", "Match routes advertised by the BGP neighbor with this exact IPv4 address")
	cmd.Flags().Var(flags.CIDRSliceFlag(), matchPrefixesFlag, "Match routes whose prefix equals one of the given IPv4 networks")

	cmd.Flags().Int32(setLocalPreferenceFlag, 0, "The BGP LOCAL_PREF to set on a permitted route. Higher values are preferred during best-path selection, the BGP default is 100. Ignored for DENY rules.")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag, filterIdFlag, actionFlag.Name())
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	minPrefixLength := flags.FlagToInt32Pointer(p, cmd, matchMinPrefixLengthFlag)
	maxPrefixLength := flags.FlagToInt32Pointer(p, cmd, matchMaxPrefixLengthFlag)
	if minPrefixLength != nil && maxPrefixLength != nil && *minPrefixLength > *maxPrefixLength {
		return nil, &errors.FlagValidationError{
			Flag:    matchMinPrefixLengthFlag,
			Details: fmt.Sprintf("must not be greater than --%s", matchMaxPrefixLengthFlag),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringValue(p, cmd, gatewayIdFlag),
		FilterId:        flags.FlagToStringValue(p, cmd, filterIdFlag),

		Action:   vpn.CreateGatewayBGPFilterRulePayloadAction(actionFlag.Get()),
		Sequence: flags.FlagToInt32Pointer(p, cmd, sequenceFlag),

		MatchAsPathContainsAny: flags.FlagToInt64SliceValue(p, cmd, matchAsPathContainsAnyFlag),
		MatchCommunities:       flags.FlagToStringSliceValue(p, cmd, matchCommunitiesFlag),
		MatchFirstAsn:          flags.FlagToInt64Pointer(p, cmd, matchFirstAsnFlag),
		MatchMaxPrefixLength:   maxPrefixLength,
		MatchMinPrefixLength:   minPrefixLength,
		MatchPeer:              flags.FlagToStringPointer(p, cmd, matchPeerFlag),
		MatchPrefixes:          flags.FlagToStringSliceValue(p, cmd, matchPrefixesFlag),

		SetLocalPreference: flags.FlagToInt32Pointer(p, cmd, setLocalPreferenceFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiCreateGatewayBGPFilterRuleRequest {
	payload := vpn.CreateGatewayBGPFilterRulePayload{
		Action:   model.Action,
		Sequence: model.Sequence,
	}

	if model.MatchAsPathContainsAny != nil || model.MatchCommunities != nil || model.MatchFirstAsn != nil ||
		model.MatchMaxPrefixLength != nil || model.MatchMinPrefixLength != nil || model.MatchPeer != nil ||
		model.MatchPrefixes != nil {
		payload.Match = &vpn.CreateGatewayBGPFilterRulePayloadMatch{
			AsPathContainsAny: model.MatchAsPathContainsAny,
			Communities:       model.MatchCommunities,
			FirstASN:          model.MatchFirstAsn,
			MaxPrefixLength:   model.MatchMaxPrefixLength,
			MinPrefixLength:   model.MatchMinPrefixLength,
			Peer:              model.MatchPeer,
			Prefixes:          model.MatchPrefixes,
		}
	}

	if model.SetLocalPreference != nil {
		payload.Set = &vpn.CreateGatewayBGPFilterRulePayloadSet{
			LocalPreference: model.SetLocalPreference,
		}
	}

	req := apiClient.DefaultAPI.CreateGatewayBGPFilterRule(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId)
	return req.CreateGatewayBGPFilterRulePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, filterLabel string, rule *vpn.BGPFilterRule) error {
	return p.OutputResult(outputFormat, rule, func() error {
		if rule == nil {
			p.Outputln("vpn bgp filter rule response is empty")
			return nil
		}
		p.Outputf("Created %s rule for BGP filter %q.\nRule ID: %s\n", rule.Action, filterLabel, utils.PtrString(rule.Id))
		return nil
	})
}
