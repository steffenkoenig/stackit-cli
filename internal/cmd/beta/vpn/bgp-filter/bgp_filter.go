package bgpfilter

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/rule"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bgp-filter",
		Short: "Provides functionality for VPN BGP filters",
		Long: `Provides functionality for VPN BGP filters.

A BGP filter holds an ordered set of rules which are used to filter the routes a VPN gateway
receives via BGP. A route is evaluated against each rule in sequence order and the first match
decides the outcome. An implicit DENY is applied after the last rule, so a route which matches no
PERMIT rule is dropped. Attaching an empty filter to a tunnel therefore denies every route.

A filter is attached to a tunnel of a VPN connection via the '--tunnel1-bgp-inbound-filter-id' and
'--tunnel2-bgp-inbound-filter-id' flags of the 'stackit beta vpn connection create' command.`,
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(rule.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
