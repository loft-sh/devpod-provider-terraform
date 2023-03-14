package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-terraform/pkg/terraform"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// CreateCmd holds the cmd flags
type CreateCmd struct{}

// NewCreateCmd defines a command
func NewCreateCmd() *cobra.Command {
	cmd := &CreateCmd{}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an instance",
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

	return createCmd
}

// Run runs the command logic
func (cmd *CreateCmd) Run(
	ctx context.Context,
	providerTerraform *terraform.TerraformProvider,
	machine *provider.Machine,
	logs log.Logger,
) error {
	err := terraform.Create(providerTerraform)
	if err != nil {
		return err
	}

	return nil
}
