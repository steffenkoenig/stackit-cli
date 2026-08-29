package describe

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

func fixtureArgValues() []string {
	return []string{testRuleId}
}

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
		RuleId:    testRuleId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *vpn.ApiGetGatewayBGPFilterRuleRequest)) vpn.ApiGetGatewayBGPFilterRuleRequest {
	request := testClient.DefaultAPI.GetGatewayBGPFilterRule(testCtx, testProjectId, testRegion, testGatewayId, testFilterId, testRuleId)
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
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no args",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid rule id",
			argValues:   []string{"not-a-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing filter id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, filterIdFlag)
			}),
			isValid: false,
		},
		{
			description: "missing project id",
			argValues:   fixtureArgValues(),
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
		expectedRequest vpn.ApiGetGatewayBGPFilterRuleRequest
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

func TestJoinInt64s(t *testing.T) {
	tests := []struct {
		description string
		values      []int64
		expected    string
	}{
		{
			description: "empty",
			values:      nil,
			expected:    "",
		},
		{
			description: "single value",
			values:      []int64{65000},
			expected:    "65000",
		},
		{
			description: "multiple values",
			values:      []int64{65000, 65001},
			expected:    "65000, 65001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := joinInt64s(tt.values, ", ")
			if got != tt.expected {
				t.Errorf("joinInt64s() = %q, want %q", got, tt.expected)
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
			name: "rule with all fields",
			rule: &vpn.BGPFilterRule{
				Action:   vpn.BGPFILTERRULEACTION_PERMIT,
				Id:       &testRuleId,
				Sequence: utils.Ptr(int32(10)),
				Match: &vpn.BGPFilterRuleMatch{
					AsPathContainsAny: []int64{65000, 65001},
					Communities:       []string{"65000:100"},
					FirstASN:          utils.Ptr(int64(65000)),
					MaxPrefixLength:   utils.Ptr(int32(24)),
					MinPrefixLength:   utils.Ptr(int32(16)),
					Peer:              utils.Ptr("192.168.0.1"),
					Prefixes:          []string{"10.0.0.0/16"},
				},
				Set: &vpn.BGPFilterRuleSet{LocalPreference: utils.Ptr(int32(200))},
			},
			wantErr: false,
		},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, tt.outputFormat, testRuleId, tt.rule); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
