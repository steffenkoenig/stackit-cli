package update

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
	testCtx           = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient        = &vpn.APIClient{DefaultAPI: &vpn.DefaultAPIService{}}
	testProjectId     = uuid.NewString()
	testGatewayId     = uuid.NewString()
	testRoutingTblId  = uuid.NewString()
	testRegion        = "eu01"
	testCurrentName   = "current-name"
	testCurrentPlanId = "p500"
)

// fixtureCurrent is the gateway as the API currently reports it.
func fixtureCurrent(mods ...func(gateway *vpn.GatewayResponse)) *vpn.GatewayResponse {
	gateway := &vpn.GatewayResponse{
		AvailabilityZones: vpn.GatewayAvailabilityZones{Tunnel1: "eu01-1", Tunnel2: "eu01-2"},
		DisplayName:       testCurrentName,
		PlanId:            testCurrentPlanId,
		RoutingType:       vpn.ROUTINGTYPE_POLICY_BASED,
		Id:                &testGatewayId,
	}
	for _, mod := range mods {
		mod(gateway)
	}
	return gateway
}

func fixtureArgValues() []string {
	return []string{testGatewayId}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		nameFlag:                  "new-name",
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
		Name:      utils.Ptr("new-name"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
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
			description: "no update flags is rejected",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, nameFlag)
			}),
			isValid: false,
		},
		{
			description: "no args",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid gateway id",
			argValues:   []string{"not-a-uuid"},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "missing project id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "labels and routing table",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[labelsFlag] = "foo=bar"
				flagValues[networkRoutingTableIdFlag] = testRoutingTblId
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Labels = &map[string]string{"foo": "bar"}
				model.NetworkRoutingTableId = &testRoutingTblId
			}),
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
		current         *vpn.GatewayResponse
		expectedPayload vpn.UpdateGatewayPayload
	}{
		{
			description: "unset attributes keep their current value",
			model:       fixtureInputModel(),
			current:     fixtureCurrent(),
			expectedPayload: vpn.UpdateGatewayPayload{
				AvailabilityZones: vpn.UpdateGatewayPayloadAvailabilityZones{Tunnel1: "eu01-1", Tunnel2: "eu01-2"},
				DisplayName:       "new-name",
				PlanId:            testCurrentPlanId,
				RoutingType:       vpn.ROUTINGTYPE_POLICY_BASED,
			},
		},
		{
			description: "routing type is carried over even though it is not settable",
			model:       fixtureInputModel(),
			current: fixtureCurrent(func(gateway *vpn.GatewayResponse) {
				gateway.RoutingType = vpn.ROUTINGTYPE_BGP_ROUTE_BASED
			}),
			expectedPayload: vpn.UpdateGatewayPayload{
				AvailabilityZones: vpn.UpdateGatewayPayloadAvailabilityZones{Tunnel1: "eu01-1", Tunnel2: "eu01-2"},
				DisplayName:       "new-name",
				PlanId:            testCurrentPlanId,
				RoutingType:       vpn.ROUTINGTYPE_BGP_ROUTE_BASED,
			},
		},
		{
			description: "existing bgp config is preserved and local asn overlaid",
			model: fixtureInputModel(func(model *inputModel) {
				model.BgpLocalAsn = utils.Ptr(int64(64513))
			}),
			current: fixtureCurrent(func(gateway *vpn.GatewayResponse) {
				gateway.Bgp = &vpn.BGPGatewayConfig{LocalAsn: 64512, OverrideAdvertisedRoutes: []string{"10.0.0.0/24"}}
			}),
			expectedPayload: vpn.UpdateGatewayPayload{
				AvailabilityZones: vpn.UpdateGatewayPayloadAvailabilityZones{Tunnel1: "eu01-1", Tunnel2: "eu01-2"},
				Bgp:               &vpn.BGPGatewayConfig{LocalAsn: 64513, OverrideAdvertisedRoutes: []string{"10.0.0.0/24"}},
				DisplayName:       "new-name",
				PlanId:            testCurrentPlanId,
				RoutingType:       vpn.ROUTINGTYPE_POLICY_BASED,
			},
		},
		{
			description: "routing table id overlaid onto existing network config",
			model: fixtureInputModel(func(model *inputModel) {
				model.NetworkRoutingTableId = &testRoutingTblId
			}),
			current: fixtureCurrent(func(gateway *vpn.GatewayResponse) {
				gateway.NetworkConfig = &vpn.NetworkConfig{PredefinedNetworkPrefix: []string{"10.0.0.0/28"}}
			}),
			expectedPayload: vpn.UpdateGatewayPayload{
				AvailabilityZones: vpn.UpdateGatewayPayloadAvailabilityZones{Tunnel1: "eu01-1", Tunnel2: "eu01-2"},
				DisplayName:       "new-name",
				NetworkConfig: &vpn.NetworkConfig{
					PredefinedNetworkPrefix: []string{"10.0.0.0/28"},
					RoutingTableId:          &testRoutingTblId,
				},
				PlanId:      testCurrentPlanId,
				RoutingType: vpn.ROUTINGTYPE_POLICY_BASED,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient, tt.current)

			expected := testClient.DefaultAPI.
				UpdateGateway(testCtx, testProjectId, testRegion, testGatewayId).
				UpdateGatewayPayload(tt.expectedPayload)

			diff := cmp.Diff(request, expected,
				cmp.AllowUnexported(expected, vpn.DefaultAPIService{}),
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
		name    string
		async   bool
		gateway *vpn.GatewayResponse
		wantErr bool
	}{
		{name: "nil gateway", wantErr: false},
		{name: "gateway set", gateway: fixtureCurrent(), wantErr: false},
		{name: "async", async: true, gateway: fixtureCurrent(), wantErr: false},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, "", tt.async, testCurrentName, tt.gateway); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
