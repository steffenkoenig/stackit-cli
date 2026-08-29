package create

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

type testCtxKey struct{}

var (
	testCtx       = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient    = &vpn.APIClient{DefaultAPI: &vpn.DefaultAPIService{}}
	testProjectId = uuid.NewString()
	testGatewayId = uuid.NewString()
	testFilterId  = uuid.NewString()
	testRuleId    = uuid.NewString()
	testRegion    = "eu01"
	testAction    = vpn.CREATEGATEWAYBGPFILTERRULEPAYLOADACTION_PERMIT
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		gatewayIdFlag:             testGatewayId,
		filterIdFlag:              testFilterId,
		actionFlag.Name():         string(testAction),
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		GatewayId: testGatewayId,
		FilterId:  testFilterId,
		Action:    testAction,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiCreateGatewayBGPFilterRuleRequest)) vpn.ApiCreateGatewayBGPFilterRuleRequest {
	request := testClient.DefaultAPI.CreateGatewayBGPFilterRule(testCtx, testProjectId, testRegion, testGatewayId, testFilterId)
	request = request.CreateGatewayBGPFilterRulePayload(vpn.CreateGatewayBGPFilterRulePayload{
		Action: testAction,
	})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description:   "base",
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "all match criteria",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[sequenceFlag] = "20"
				flagValues[matchAsPathContainsAnyFlag] = "65000,65001"
				flagValues[matchCommunitiesFlag] = "65000:100"
				flagValues[matchFirstAsnFlag] = "65000"
				flagValues[matchMinPrefixLengthFlag] = "16"
				flagValues[matchMaxPrefixLengthFlag] = "24"
				flagValues[matchPeerFlag] = "192.168.0.1"
				flagValues[matchPrefixesFlag] = "10.0.0.0/16"
				flagValues[setLocalPreferenceFlag] = "200"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Sequence = utils.Ptr(int32(20))
				model.MatchAsPathContainsAny = []int64{65000, 65001}
				model.MatchCommunities = []string{"65000:100"}
				model.MatchFirstAsn = utils.Ptr(int64(65000))
				model.MatchMinPrefixLength = utils.Ptr(int32(16))
				model.MatchMaxPrefixLength = utils.Ptr(int32(24))
				model.MatchPeer = utils.Ptr("192.168.0.1")
				model.MatchPrefixes = []string{"10.0.0.0/16"}
				model.SetLocalPreference = utils.Ptr(int32(200))
			}),
		},
		{
			description: "min prefix length greater than max prefix length",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[matchMinPrefixLengthFlag] = "24"
				flagValues[matchMaxPrefixLengthFlag] = "16"
			}),
			isValid: false,
		},
		{
			description: "invalid match prefix",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[matchPrefixesFlag] = "not-a-cidr"
			}),
			isValid: false,
		},
		{
			description: "invalid action",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[actionFlag.Name()] = "ALLOW"
			}),
			isValid: false,
		},
		{
			description: "missing action",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, actionFlag.Name())
			}),
			isValid: false,
		},
		{
			description: "missing filter id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, filterIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing project id",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		expectedRequest vpn.ApiCreateGatewayBGPFilterRuleRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
		},
		{
			description: "with match and set",
			model: fixtureInputModel(func(model *inputModel) {
				model.Sequence = utils.Ptr(int32(20))
				model.MatchPrefixes = []string{"10.0.0.0/16"}
				model.MatchPeer = utils.Ptr("192.168.0.1")
				model.SetLocalPreference = utils.Ptr(int32(200))
			}),
			expectedRequest: fixtureRequest(func(request *vpn.ApiCreateGatewayBGPFilterRuleRequest) {
				*request = request.CreateGatewayBGPFilterRulePayload(vpn.CreateGatewayBGPFilterRulePayload{
					Action:   testAction,
					Sequence: utils.Ptr(int32(20)),
					Match: &vpn.CreateGatewayBGPFilterRulePayloadMatch{
						Prefixes: []string{"10.0.0.0/16"},
						Peer:     utils.Ptr("192.168.0.1"),
					},
					Set: &vpn.CreateGatewayBGPFilterRulePayloadSet{
						LocalPreference: utils.Ptr(int32(200)),
					},
				})
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient)

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest, vpn.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name         string
		outputFormat string
		rule         *vpn.BGPFilterRule
		wantErr      bool
	}{
		{
			name:    "nil rule",
			wantErr: false,
		},
		{
			name:    "empty rule",
			rule:    &vpn.BGPFilterRule{},
			wantErr: false,
		},
		{
			name:    "rule set",
			rule:    &vpn.BGPFilterRule{Action: vpn.BGPFILTERRULEACTION_PERMIT, Id: &testRuleId},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, "my-filter", tt.rule); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
