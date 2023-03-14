package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/loft-sh/devpod-provider-terraform/pkg/terraform"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/provider"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
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

	// get private key
	privateKey, err := ssh.GetPrivateKeyRawBase(providerTerraform.Config.MachineFolder)

	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	// get external address
	externalIP, err := terraform.GetExternalIP(providerTerraform)
	if err != nil || externalIP == "" {
		return fmt.Errorf(
			"instance %s doesn't have an external nat ip",
			providerTerraform.Config.MachineID,
		)
	}

	sshClient, err := ssh.NewSSHClient("devpod", externalIP+":22", privateKey)

	if err != nil {
		return errors.Wrap(err, "create ssh client")
	}

	defer sshClient.Close()

	// run command
	return ssh.Run(sshClient, command, os.Stdin, os.Stdout, os.Stderr)
}
