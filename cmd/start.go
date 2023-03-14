package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-terraform/pkg/terraform"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// StartCmd holds the cmd flags
type StartCmd struct{}

// NewStartCmd defines a command
func NewStartCmd() *cobra.Command {
	cmd := &StartCmd{}
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an instance",
		RunE: func(_ *cobra.Command, args []string) error {
			terraformProvider, err := terraform.NewProvider(log.Default)
			if err != nil {
				return err
			}

			return cmd.Run(
				context.Background(),
				terraformProvider,
				provider.FromEnvironment(),
				log.Default,
			)
		},
	}

	return startCmd
}

// Run runs the command logic
func (cmd *StartCmd) Run(
	ctx context.Context,
	providerTerraform *terraform.TerraformProvider,
	machine *provider.Machine,
	logs log.Logger,
) error {

	return nil
}
