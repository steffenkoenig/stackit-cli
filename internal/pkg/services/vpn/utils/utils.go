package utils

import (
	"context"
	"fmt"
	"slices"
	"strings"

	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

func GetGatewayName(ctx context.Context, client vpn.DefaultAPI, projectId, region, gatewayId string) (string, error) {
	resp, err := client.GetGateway(ctx, projectId, region, gatewayId).Execute()
	if err != nil {
		return "", fmt.Errorf("get gateway: %w", err)
	}
	if resp != nil {
		return resp.DisplayName, nil
	}
	return "", nil
}

// GetBGPFilterName returns the display name of the BGP filter, or an empty string if it cannot be determined.
func GetBGPFilterName(ctx context.Context, client vpn.DefaultAPI, projectId, region, gatewayId, filterId string) (string, error) {
	resp, err := client.GetGatewayBGPFilter(ctx, projectId, region, gatewayId, filterId).Execute()
	if err != nil {
		return "", fmt.Errorf("get bgp filter: %w", err)
	}
	if resp != nil {
		return resp.DisplayName, nil
	}
	return "", nil
}

// deprecatedIntegrityAlgorithms lists the integrity algorithms which are still accepted by the API
// but should no longer be used.
var deprecatedIntegrityAlgorithms = []vpn.PhaseIntegrityAlgorithmsInner{
	vpn.PHASEINTEGRITYALGORITHMSINNER_SHA1,
}

// recommendedIntegrityAlgorithms lists the integrity algorithms users should migrate to.
var recommendedIntegrityAlgorithms = []vpn.PhaseIntegrityAlgorithmsInner{
	vpn.PHASEINTEGRITYALGORITHMSINNER_SHA2_256,
	vpn.PHASEINTEGRITYALGORITHMSINNER_SHA2_384,
	vpn.PHASEINTEGRITYALGORITHMSINNER_SHA2_512,
}

// WarnDeprecatedIntegrityAlgorithms prints a warning for every deprecated integrity algorithm
// configured via the given flag, pointing to the recommended alternatives.
func WarnDeprecatedIntegrityAlgorithms(p *print.Printer, flagName string, algorithms []vpn.PhaseIntegrityAlgorithmsInner) {
	recommended := make([]string, 0, len(recommendedIntegrityAlgorithms))
	for _, algorithm := range recommendedIntegrityAlgorithms {
		recommended = append(recommended, string(algorithm))
	}

	for _, algorithm := range algorithms {
		if !slices.Contains(deprecatedIntegrityAlgorithms, algorithm) {
			continue
		}
		p.Warn(
			"The integrity algorithm %q configured via --%s is deprecated and should no longer be used. Please use one of %s instead.\n",
			string(algorithm),
			flagName,
			strings.Join(recommended, ", "),
		)
	}
}
