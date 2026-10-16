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
	ruleIdArg = "RULE_ID"

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
	vpn.AllowedUpdateGatewayBGPFilterRulePayloadActionEnumValues,
	"The action to take if a route matches all criteria of the rule",
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId string
	FilterId  string
	RuleId    string

	Action   vpn.UpdateGatewayBGPFilterRulePayloadAction
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
		Use:   fmt.Sprintf("update %s", ruleIdArg),
		Short: "Updates a VPN BGP filter rule",
		Long: `Updates a VPN BGP filter rule.

The rule is replaced as a whole, so every match criterion which should be kept has to be passed
again. Criteria which are not passed are removed from the rule.`,
		Args: args.SingleArg(ruleIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Turn the rule with the ID "zzz" into a rule which permits every route`,
				"$ stackit beta vpn bgp-filter rule update zzz --gateway-id xxx --filter-id yyy --action PERMIT --sequence 10",
			),
			examples.NewExample(
				`Change the rule with the ID "zzz" to deny the routes 10.0.0.0/16 and 192.168.0.0/16`,
				"$ stackit beta vpn bgp-filter rule update zzz --gateway-id xxx --filter-id yyy --action DENY --sequence 10 --match-prefixes 10.0.0.0/16,192.168.0.0/16",
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

			prompt := fmt.Sprintf("Are you sure you want to update the rule %q of the BGP filter %q?", model.RuleId, filterLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpn bgp filter rule: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.RuleId, filterLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().Var(flags.UUIDFlag(), filterIdFlag, "BGP filter ID")
	actionFlag.Register(cmd.Flags())
	cmd.Flags().Int32(sequenceFlag, 0, "The evaluation order of the rule. Lower numbers are evaluated first and the number must be unique within the filter.")

	cmd.Flags().Int64Slice(matchAsPathContainsAnyFlag, nil, "Match routes whose AS-PATH contains any of the given ASNs")
	cmd.Flags().StringSlice(matchCommunitiesFlag, nil, "Match routes carrying any of the given BGP standard communities in the 'asn:value' format (RFC 1997)")
	cmd.Flags().Int64(matchFirstAsnFlag, 0, "Match routes whose first ASN in the AS-PATH, i.e. the immediate neighbor, equals the given ASN")
	cmd.Flags().Int32(matchMaxPrefixLengthFlag, 0, "Match routes with a subnet mask length of at most this value")
	cmd.Flags().Int32(matchMinPrefixLengthFlag, 0, "Match routes with a subnet mask length of at least this value")
	cmd.Flags().String(matchPeerFlag, "", "Match routes advertised by the BGP neighbor with this exact IPv4 address")
	cmd.Flags().Var(flags.CIDRSliceFlag(), matchPrefixesFlag, "Match routes whose prefix equals one of the given IPv4 networks")

	cmd.Flags().Int32(setLocalPreferenceFlag, 0, "The BGP LOCAL_PREF to set on a permitted route. Higher values are preferred during best-path selection, the BGP default is 100. Ignored for DENY rules.")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag, filterIdFlag, actionFlag.Name(), sequenceFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
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
		RuleId:          inputArgs[0],

		Action:   vpn.UpdateGatewayBGPFilterRulePayloadAction(actionFlag.Get()),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) vpn.ApiUpdateGatewayBGPFilterRuleRequest {
	payload := vpn.UpdateGatewayBGPFilterRulePayload{
		Action:   model.Action,
		Sequence: model.Sequence,
	}

	if model.MatchAsPathContainsAny != nil || model.MatchCommunities != nil || model.MatchFirstAsn != nil ||
		model.MatchMaxPrefixLength != nil || model.MatchMinPrefixLength != nil || model.MatchPeer != nil ||
		model.MatchPrefixes != nil {
		payload.Match = &vpn.UpdateGatewayBGPFilterRulePayloadMatch{
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
		payload.Set = &vpn.UpdateGatewayBGPFilterRulePayloadSet{
			LocalPreference: model.SetLocalPreference,
		}
	}

	req := apiClient.DefaultAPI.UpdateGatewayBGPFilterRule(ctx, model.ProjectId, model.Region, model.GatewayId, model.FilterId, model.RuleId)
	return req.UpdateGatewayBGPFilterRulePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, ruleId, filterLabel string, rule *vpn.BGPFilterRule) error {
	return p.OutputResult(outputFormat, rule, func() error {
		if rule == nil {
			p.Outputln("vpn bgp filter rule response is empty")
			return nil
		}
		p.Outputf("Updated rule %q of BGP filter %q.\n", ruleId, filterLabel)
		return nil
	})
}
