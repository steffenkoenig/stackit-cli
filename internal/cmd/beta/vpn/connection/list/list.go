package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	vpn "github.com/stackitcloud/stackit-sdk-go/services/vpn/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/vpn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	gatewayIdFlag = "gateway-id"
	limitFlag     = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	GatewayId *string
	Limit     *int64
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all VPN connections of a gateway",
		Long:  "Lists all VPN connections of a gateway.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all VPN connections of a gateway`,
				"$ stackit beta vpn connection list --gateway-id xxx"),
			examples.NewExample(
				`List up to 4 VPN connections of a gateway`,
				"$ stackit beta vpn connection list --gateway-id xxx --limit 4"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list VPN connections: %w", err)
			}

			// Truncate output
			if resp != nil && model.Limit != nil && len(resp.Connections) > int(*model.Limit) {
				resp.Connections = resp.Connections[:*model.Limit]
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), gatewayIdFlag, "Gateway ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, gatewayIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		GatewayId:       flags.FlagToStringPointer(p, cmd, gatewayIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *vpn.APIClient) (vpn.ApiListGatewayConnectionsRequest, error) {
	req := apiClient.DefaultAPI.ListGatewayConnections(ctx, model.ProjectId, model.Region, *model.GatewayId)
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *vpn.ConnectionList) error {
	if resp == nil || resp.Connections == nil {
		return fmt.Errorf("list connections response is empty")
	}

	return p.OutputResult(model.OutputFormat, resp.Connections, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "ENABLED", "LABELS")
		for _, c := range resp.Connections {
			id := utils.PtrString(c.Id)
			name := c.DisplayName
			enabled := utils.PtrString(c.Enabled)
			var labels string
			if c.Labels != nil {
				labels = utils.JoinStringMap(*c.Labels, "=", ", ")
			}
			table.AddRow(id, name, enabled, labels)
		}
		p.Outputln(table.Render())
		return nil
	})
}
