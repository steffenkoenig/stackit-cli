package rule

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/rule/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/rule/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/rule/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/rule/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/vpn/bgp-filter/rule/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Provides functionality for VPN BGP filter rules",
		Long: `Provides functionality for VPN BGP filter rules.

All non-empty match criteria of a rule are AND-combined. Rules within a filter are evaluated in
sequence order and the first matching rule decides the outcome. An implicit DENY follows the last
rule. A filter holds at most 10 rules.`,
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
	cmd.AddCommand(update.NewCmd(params))
}
