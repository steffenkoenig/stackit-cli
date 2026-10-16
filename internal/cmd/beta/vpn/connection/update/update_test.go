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
	testCtx          = context.WithValue(context.Background(), testCtxKey{}, "foo")
	testClient       = &vpn.APIClient{DefaultAPI: &vpn.DefaultAPIService{}}
	testProjectId    = uuid.NewString()
	testGatewayId    = uuid.NewString()
	testConnectionId = uuid.NewString()
	testFilterId     = uuid.NewString()
	testRegion       = "eu01"
	testCurrentName  = "current-connection"
)

func fixtureTunnel(mods ...func(tunnel *vpn.TunnelConfiguration)) vpn.TunnelConfiguration {
	tunnel := vpn.TunnelConfiguration{
		RemoteAddress: "1.2.3.4",
		Phase1: vpn.TunnelConfigurationPhase1{
			EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
			IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
		},
		Phase2: vpn.TunnelConfigurationPhase2{
			EncryptionAlgorithms: []vpn.PhaseEncryptionAlgorithmsInner{"aes256"},
			IntegrityAlgorithms:  []vpn.PhaseIntegrityAlgorithmsInner{"sha2_256"},
		},
	}
	for _, mod := range mods {
		mod(&tunnel)
	}
	return tunnel
}

// fixtureCurrent is the connection as the API currently reports it.
func fixtureCurrent(mods ...func(conn *vpn.ConnectionResponse)) *vpn.ConnectionResponse {
	conn := &vpn.ConnectionResponse{
		DisplayName: testCurrentName,
		Enabled:     utils.Ptr(true),
		Id:          &testConnectionId,
		Tunnel1:     fixtureTunnel(),
		Tunnel2:     fixtureTunnel(func(t *vpn.TunnelConfiguration) { t.RemoteAddress = "5.6.7.8" }),
	}
	for _, mod := range mods {
		mod(conn)
	}
	return conn
}

func fixtureArgValues() []string {
	return []string{testConnectionId}
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		gatewayIdFlag:             testGatewayId,
		displayNameFlag:           "new-name",
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
		GatewayId:    testGatewayId,
		ConnectionId: testConnectionId,
		DisplayName:  utils.Ptr("new-name"),
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
				delete(flagValues, displayNameFlag)
			}),
			isValid: false,
		},
		{
			description: "missing gateway id",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, gatewayIdFlag)
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
		{
			description: "attach inbound filter to tunnel 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				delete(flagValues, displayNameFlag)
				flagValues[tunnel1BgpInboundFilterIdFlag] = testFilterId
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Tunnel1.BgpInboundFilterId = &testFilterId
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			testutils.TestParseInput(t, NewCmd, parseInput, tt.expectedModel, tt.argValues, tt.flagValues, tt.isValid)
		})
	}
}

func TestValidateBgpInboundFilters(t *testing.T) {
	tests := []struct {
		description string
		model       *inputModel
		current     *vpn.ConnectionResponse
		wantErr     bool
	}{
		{
			description: "no filter set",
			model:       fixtureInputModel(),
			current:     fixtureCurrent(),
			wantErr:     false,
		},
		{
			description: "filter with no bgp session anywhere is rejected",
			model: fixtureInputModel(func(model *inputModel) {
				model.Tunnel1.BgpInboundFilterId = &testFilterId
			}),
			current: fixtureCurrent(),
			wantErr: true,
		},
		{
			description: "filter is allowed when the asn comes from a flag",
			model: fixtureInputModel(func(model *inputModel) {
				model.Tunnel1.BgpInboundFilterId = &testFilterId
				model.Tunnel1.BgpRemoteAsn = utils.Ptr(int64(65000))
			}),
			current: fixtureCurrent(),
			wantErr: false,
		},
		{
			description: "filter is allowed when the tunnel already has a bgp session",
			model: fixtureInputModel(func(model *inputModel) {
				model.Tunnel1.BgpInboundFilterId = &testFilterId
			}),
			current: fixtureCurrent(func(conn *vpn.ConnectionResponse) {
				conn.Tunnel1.Bgp = &vpn.BGPTunnelConfig{RemoteAsn: 65000}
			}),
			wantErr: false,
		},
		{
			description: "tunnel 2 is checked as well",
			model: fixtureInputModel(func(model *inputModel) {
				model.Tunnel2.BgpInboundFilterId = &testFilterId
			}),
			current: fixtureCurrent(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := validateBgpInboundFilters(tt.model, tt.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBgpInboundFilters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description     string
		model           *inputModel
		current         *vpn.ConnectionResponse
		expectedPayload vpn.UpdateGatewayConnectionPayload
	}{
		{
			description: "unset attributes keep their current value",
			model:       fixtureInputModel(),
			current:     fixtureCurrent(),
			expectedPayload: vpn.UpdateGatewayConnectionPayload{
				DisplayName: "new-name",
				Enabled:     utils.Ptr(true),
				Tunnel1:     fixtureTunnel(),
				Tunnel2:     fixtureTunnel(func(t *vpn.TunnelConfiguration) { t.RemoteAddress = "5.6.7.8" }),
			},
		},
		{
			description: "attaching a filter preserves the existing remote asn",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Tunnel1.BgpInboundFilterId = &testFilterId
			}),
			current: fixtureCurrent(func(conn *vpn.ConnectionResponse) {
				conn.Tunnel1.Bgp = &vpn.BGPTunnelConfig{RemoteAsn: 65000}
			}),
			expectedPayload: vpn.UpdateGatewayConnectionPayload{
				DisplayName: testCurrentName,
				Enabled:     utils.Ptr(true),
				Tunnel1: fixtureTunnel(func(t *vpn.TunnelConfiguration) {
					t.Bgp = &vpn.BGPTunnelConfig{
						RemoteAsn:       65000,
						InboundFilterId: *vpn.NewNullableString(&testFilterId),
					}
				}),
				Tunnel2: fixtureTunnel(func(t *vpn.TunnelConfiguration) { t.RemoteAddress = "5.6.7.8" }),
			},
		},
		{
			description: "phase config is overlaid onto the current tunnel",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Tunnel1.Phase1IntegrityAlgorithms = []vpn.PhaseIntegrityAlgorithmsInner{"sha2_512"}
			}),
			current: fixtureCurrent(),
			expectedPayload: vpn.UpdateGatewayConnectionPayload{
				DisplayName: testCurrentName,
				Enabled:     utils.Ptr(true),
				Tunnel1: fixtureTunnel(func(t *vpn.TunnelConfiguration) {
					t.Phase1.IntegrityAlgorithms = []vpn.PhaseIntegrityAlgorithmsInner{"sha2_512"}
				}),
				Tunnel2: fixtureTunnel(func(t *vpn.TunnelConfiguration) { t.RemoteAddress = "5.6.7.8" }),
			},
		},
		{
			description: "disabling the connection keeps everything else",
			model: fixtureInputModel(func(model *inputModel) {
				model.DisplayName = nil
				model.Enabled = utils.Ptr(false)
			}),
			current: fixtureCurrent(),
			expectedPayload: vpn.UpdateGatewayConnectionPayload{
				DisplayName: testCurrentName,
				Enabled:     utils.Ptr(false),
				Tunnel1:     fixtureTunnel(),
				Tunnel2:     fixtureTunnel(func(t *vpn.TunnelConfiguration) { t.RemoteAddress = "5.6.7.8" }),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request := buildRequest(testCtx, tt.model, testClient, tt.current)

			expected := testClient.DefaultAPI.
				UpdateGatewayConnection(testCtx, testProjectId, testRegion, testGatewayId, testConnectionId).
				UpdateGatewayConnectionPayload(tt.expectedPayload)

			diff := cmp.Diff(request, expected,
				cmp.AllowUnexported(expected, vpn.DefaultAPIService{}),
				cmpopts.EquateComparable(testCtx),
				cmp.AllowUnexported(vpn.NullableString{}),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

// The API never returns the pre-shared key, so it must only be sent when passed explicitly.
func TestPreSharedKeyOnlySentWhenGiven(t *testing.T) {
	current := fixtureCurrent()

	notGiven := mergeTunnelConfiguration(&tunnelInputModel{}, &current.Tunnel1)
	if notGiven.PreSharedKey != nil {
		t.Errorf("expected no pre-shared key to be sent, got %v", *notGiven.PreSharedKey)
	}

	given := mergeTunnelConfiguration(&tunnelInputModel{PreSharedKey: utils.Ptr("new-psk")}, &current.Tunnel1)
	if given.PreSharedKey == nil || *given.PreSharedKey != "new-psk" {
		t.Errorf("expected the given pre-shared key to be sent, got %v", given.PreSharedKey)
	}
}

func TestOutputResult(t *testing.T) {
	tests := []struct {
		name    string
		resp    *vpn.ConnectionResponse
		wantErr bool
	}{
		{name: "nil response", wantErr: false},
		{name: "response set", resp: fixtureCurrent(), wantErr: false},
	}
	params := testparams.NewTestParams()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(params.Printer, "", testCurrentName, tt.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
