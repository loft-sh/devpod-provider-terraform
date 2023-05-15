package cmd

import (
	"context"

	"github.com/loft-sh/devpod-provider-terraform/pkg/options"
	"github.com/loft-sh/devpod-provider-terraform/pkg/terraform"

	"github.com/loft-sh/devpod/pkg/config"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// InitCmd holds the cmd flags
type InitCmd struct{}

// NewInitCmd defines a init
func NewInitCmd() *cobra.Command {
	cmd := &InitCmd{}
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init account",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run(
				context.Background(),
				provider.FromEnvironment(),
				log.Default,
			)
		},
	}

	return initCmd
}

// Run runs the init logic
func (cmd *InitCmd) Run(
	ctx context.Context,
	machine *provider.Machine,
	logs log.Logger,
) error {
	devpodPath, err := config.GetConfigDir()
	if err != nil {
		return err
	}

	terraformPath := devpodPath + "/bin/terraform"

	project, err := options.FromEnvOrError(options.TERRAFORM_PROJECT)
	if err != nil {
		return err
	}

	// create provider
	provider := &terraform.TerraformProvider{
		Log:     logs,
		Bin:     terraformPath,
		Project: project,
	}

	err = terraform.Install(provider)
	if err != nil {
		return err
	}

	return nil
}
