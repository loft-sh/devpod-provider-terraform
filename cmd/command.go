package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/loft-sh/devpod-provider-terraform/pkg/terraform"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/spf13/cobra"
)

// CommandCmd holds the cmd flags
type CommandCmd struct{}

// NewCommandCmd defines a command
func NewCommandCmd() *cobra.Command {
	cmd := &CommandCmd{}
	commandCmd := &cobra.Command{
		Use:   "command",
		Short: "Command an instance",
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

	return commandCmd
}

// Run runs the command logic
func (cmd *CommandCmd) Run(
	ctx context.Context,
	providerTerraform *terraform.TerraformProvider,
	machine *provider.Machine,
	logs log.Logger,
) error {
	command := os.Getenv("COMMAND")

	if command == "" {
		return fmt.Errorf("command environment variable is missing")
	}

	return terraform.Command(providerTerraform, command)
}
