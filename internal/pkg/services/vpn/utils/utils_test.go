package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/testparams"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	testGatewayName = "test-gateway-01"
	testRegion      = "eu01"
)

var (
	testProjectId = uuid.NewString()
	testGatewayId = uuid.NewString()
)

type mockSettings struct {
	getGatewayFails bool
	getGatewayResp  *vpn.GatewayResponse
}

func newAPIMock(settings *mockSettings) vpn.DefaultAPI {
	return &vpn.DefaultAPIServiceMock{
		GetGatewayExecuteMock: utils.Ptr(func(_ vpn.ApiGetGatewayRequest) (*vpn.GatewayResponse, error) {
			if settings.getGatewayFails {
				return nil, fmt.Errorf("could not get gateway details")
			}

			return settings.getGatewayResp, nil
		}),
	}
}

func TestGetGatewayName(t *testing.T) {
	tests := []struct {
		description     string
		getGatewayResp  *vpn.GatewayResponse
		getGatewayFails bool
		isValid         bool
		expectedOutput  string
	}{
		{
			description: "base",
			getGatewayResp: &vpn.GatewayResponse{
				DisplayName: testGatewayName,
			},
			isValid:        true,
			expectedOutput: testGatewayName,
		},
		{
			description:     "get gateway fails",
			getGatewayFails: true,
			isValid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := newAPIMock(&mockSettings{
				getGatewayFails: tt.getGatewayFails,
				getGatewayResp:  tt.getGatewayResp,
			})

			output, err := GetGatewayName(context.Background(), client, testProjectId, testRegion, testGatewayId)

			if tt.isValid && err != nil {
				t.Errorf("failed on valid input")
			}
			if !tt.isValid && err == nil {
				t.Errorf("did not fail on invalid input")
			}
			if !tt.isValid {
				return
			}
			if output != tt.expectedOutput {
				t.Errorf("expected output to be %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

func TestWarnDeprecatedIntegrityAlgorithms(t *testing.T) {
	tests := []struct {
		description   string
		algorithms    []vpn.PhaseIntegrityAlgorithmsInner
		expectWarning bool
	}{
		{
			description:   "no algorithms",
			algorithms:    nil,
			expectWarning: false,
		},
		{
			description:   "sha2_256 only",
			algorithms:    []vpn.PhaseIntegrityAlgorithmsInner{vpn.PHASEINTEGRITYALGORITHMSINNER_SHA2_256},
			expectWarning: false,
		},
		{
			description:   "sha2_512 only",
			algorithms:    []vpn.PhaseIntegrityAlgorithmsInner{vpn.PHASEINTEGRITYALGORITHMSINNER_SHA2_512},
			expectWarning: false,
		},
		{
			description:   "sha1",
			algorithms:    []vpn.PhaseIntegrityAlgorithmsInner{vpn.PHASEINTEGRITYALGORITHMSINNER_SHA1},
			expectWarning: true,
		},
		{
			description: "sha1 among others",
			algorithms: []vpn.PhaseIntegrityAlgorithmsInner{
				vpn.PHASEINTEGRITYALGORITHMSINNER_SHA2_384,
				vpn.PHASEINTEGRITYALGORITHMSINNER_SHA1,
			},
			expectWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			params := testparams.NewTestParams()
			WarnDeprecatedIntegrityAlgorithms(params.Printer, "tunnel1-phase1-integrity-algorithms", tt.algorithms)

			warning := params.Err.String()
			if tt.expectWarning {
				if !strings.Contains(warning, "sha1") {
					t.Errorf("expected a warning mentioning sha1, got %q", warning)
				}
				for _, recommended := range []string{"sha2_256", "sha2_384", "sha2_512"} {
					if !strings.Contains(warning, recommended) {
						t.Errorf("expected the warning to recommend %q, got %q", recommended, warning)
					}
				}
				if !strings.Contains(warning, "tunnel1-phase1-integrity-algorithms") {
					t.Errorf("expected the warning to name the flag, got %q", warning)
				}
				if !strings.Contains(warning, deprecatedIntegrityAlgorithmRemovalDate) {
					t.Errorf("expected the warning to name the removal date, got %q", warning)
				}
			} else if warning != "" {
				t.Errorf("expected no warning, got %q", warning)
			}
		})
	}
}
