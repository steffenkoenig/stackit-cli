package list

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
)

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		gatewayIdFlag:             testGatewayId,
		filterIdFlag:              testFilterId,
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
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiListGatewayBGPFilterRulesRequest)) vpn.ApiListGatewayBGPFilterRulesRequest {
	request := testClient.DefaultAPI.ListGatewayBGPFilterRules(testCtx, testProjectId, testRegion, testGatewayId, testFilterId)
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
			description: "with limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "4"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Limit = utils.Ptr(int64(4))
			}),
		},
		{
			description: "invalid limit",
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[limitFlag] = "0"
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
		expectedRequest vpn.ApiListGatewayBGPFilterRulesRequest
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			expectedRequest: fixtureRequest(),
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

func TestSummarizeMatch(t *testing.T) {
	tests := []struct {
		description string
		match       *vpn.BGPFilterRuleMatch
		expected    string
	}{
		{
			description: "nil match",
			match:       nil,
			expected:    "any",
		},
		{
			description: "empty match",
			match:       &vpn.BGPFilterRuleMatch{},
			expected:    "any",
		},
		{
			description: "prefixes and peer",
			match: &vpn.BGPFilterRuleMatch{
				Prefixes: []string{"10.0.0.0/16", "192.168.0.0/16"},
				Peer:     utils.Ptr("192.168.0.1"),
			},
			expected: "prefixes: 10.0.0.0/16, 192.168.0.0/16\npeer: 192.168.0.1",
		},
		{
			description: "as path and communities",
			match: &vpn.BGPFilterRuleMatch{
				AsPathContainsAny: []int64{65000, 65001},
				Communities:       []string{"65000:100"},
				FirstASN:          utils.Ptr(int64(65000)),
			},
			expected: "first ASN: 65000\nAS-PATH contains any: 65000, 65001\ncommunities: 65000:100",
		},
		{
			description: "prefix lengths",
			match: &vpn.BGPFilterRuleMatch{
				MinPrefixLength: utils.Ptr(int32(16)),
				MaxPrefixLength: utils.Ptr(int32(24)),
			},
			expected: "min prefix length: 16\nmax prefix length: 24",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := summarizeMatch(tt.match)
			if got != tt.expected {
				t.Errorf("summarizeMatch() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name         string
		outputFormat string
		rules        []vpn.BGPFilterRule
		wantErr      bool
	}{
		{
			name:    "no rules",
			wantErr: false,
		},
		{
			name:    "empty rule",
			rules:   []vpn.BGPFilterRule{{}},
			wantErr: false,
		},
		{
			name: "rules set",
			rules: []vpn.BGPFilterRule{{
				Action:   vpn.BGPFILTERRULEACTION_PERMIT,
				Id:       &testRuleId,
				Sequence: utils.Ptr(int32(10)),
				Match:    &vpn.BGPFilterRuleMatch{Prefixes: []string{"10.0.0.0/16"}},
				Set:      &vpn.BGPFilterRuleSet{LocalPreference: utils.Ptr(int32(200))},
			}},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, testFilterId, tt.rules); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
