package terraform

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/loft-sh/devpod-provider-terraform/pkg/options"
	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/config"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"

	cp "github.com/otiai10/copy"
)

func NewProvider(logs log.Logger) (*TerraformProvider, error) {
	providerConfig, err := options.FromEnv()
	if err != nil {
		return nil, err
	}

	devpodPath, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}

	terraformPath := devpodPath + "/bin/terraform"

	project, err := options.FromEnvOrError(options.TERRAFORM_PROJECT)
	if err != nil {
		return nil, err
	}

	// create provider
	provider := &TerraformProvider{
		Config:     providerConfig,
		Log:        logs,
		Bin:        terraformPath,
		Project:    project,
		State:      providerConfig.MachineFolder + "/main.tfstate",
		WorkingDir: providerConfig.MachineFolder + "/.terraform",
	}

	return provider, nil
}

type TerraformProvider struct {
	Config     *options.Options
	Log        log.Logger
	Bin        string
	Project    string
	State      string
	WorkingDir string
}

func EnsureProject(providerTerraform *TerraformProvider) error {
	// if project is already in place, exit
	_, err := os.Stat(providerTerraform.Config.MachineFolder + "/.terraform")
	if err == nil {
		return nil
	}

	// if project is an url, try to clone it
	if strings.Contains(providerTerraform.Project, "http://") ||
		strings.Contains(providerTerraform.Project, "https://") {
		cmd := exec.Command(
			"git",
			"clone",
			providerTerraform.Project,
			providerTerraform.Config.MachineFolder+"/.terraform",
		)
		return cmd.Run()
	}

	// else we have a path, let's copy it to destination
	_, err = os.Stat(providerTerraform.Project)
	if err != nil {
		return errors.Errorf("terraform project not found")
	}

	err = cp.Copy(providerTerraform.Project,
		providerTerraform.Config.MachineFolder+"/.terraform")
	if err != nil {
		return err
	}

	return nil
}

func Init(providerTerraform *TerraformProvider) (*tfexec.Terraform, error) {
	err := EnsureProject(providerTerraform)
	if err != nil {
		return nil, err
	}

	workingDir := providerTerraform.Config.MachineFolder + "/.terraform"
	tf, err := tfexec.NewTerraform(workingDir, providerTerraform.Bin)
	if err != nil {
		return nil, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return nil, err
	}

	return tf, nil
}

func Install(providerTerraform *TerraformProvider) error {
	err := exec.Command(providerTerraform.Bin).Run()
	if err == nil {
		return nil
	}

	destPath := filepath.Dir(providerTerraform.Bin)

	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return err
	}

	installer := &releases.ExactVersion{
		InstallDir: destPath,
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.4.0")),
	}

	_, err = installer.Install(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func Delete(providerTerraform *TerraformProvider) error {
	tf, err := Init(providerTerraform)
	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background(),
		tfexec.Lock(false),
		tfexec.Refresh(true),
		tfexec.Parallelism(99),
		tfexec.State(providerTerraform.State),
	)
	if err != nil {
		return err
	}

	return nil
}

func Command(providerTerraform *TerraformProvider, command string) error {
	// get private key
	privateKey, err := ssh.GetPrivateKeyRawBase(providerTerraform.Config.MachineFolder)

	if err != nil {
		return fmt.Errorf("load private key: %w", err)
	}

	// get external address
	externalIP, err := getExternalIP(providerTerraform)
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
	return ssh.Run(context.Background(), sshClient, command, os.Stdin, os.Stdout, os.Stderr)
}

func Create(providerTerraform *TerraformProvider) error {
	tf, err := Init(providerTerraform)
	if err != nil {
		return err
	}

	publicKeyBase, err := ssh.GetPublicKeyBase(providerTerraform.Config.MachineFolder)
	if err != nil {
		return err
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return err
	}

	err = tf.Apply(context.Background(),
		tfexec.Lock(false),
		tfexec.Refresh(true),
		tfexec.Parallelism(99),
		tfexec.State(providerTerraform.State),
		tfexec.Var("disk_image="+providerTerraform.Config.DiskImage),
		tfexec.Var("disk_size="+providerTerraform.Config.DiskSizeGB),
		tfexec.Var("instance_type="+providerTerraform.Config.MachineType),
		tfexec.Var("machine_name="+providerTerraform.Config.MachineID),
		tfexec.Var("region="+providerTerraform.Config.Zone),
		tfexec.Var("ssh_key="+string(publicKey)),
	)
	if err != nil {
		return err
	}
	err = tf.Refresh(context.Background(),
		tfexec.Lock(false),
		tfexec.State(providerTerraform.State),
		tfexec.Var("disk_image="+providerTerraform.Config.DiskImage),
		tfexec.Var("disk_size="+providerTerraform.Config.DiskSizeGB),
		tfexec.Var("instance_type="+providerTerraform.Config.MachineType),
		tfexec.Var("machine_name="+providerTerraform.Config.MachineID),
		tfexec.Var("region="+providerTerraform.Config.Zone),
		tfexec.Var("ssh_key="+string(publicKey)),
	)
	if err != nil {
		return err
	}

	return nil
}

func getExternalIP(providerTerraform *TerraformProvider) (string, error) {
	tf, err := Init(providerTerraform)
	if err != nil {
		return "", err
	}

	output, err := tf.Output(context.Background(),
		tfexec.State(providerTerraform.State),
	)
	if err != nil {
		return "", err
	}

	if output["public_ip"].Value == nil {
		return "", errors.Errorf("output not found")
	}

	return strings.ReplaceAll(string(output["public_ip"].Value), "\"", ""), nil
}

func Status(providerTerraform *TerraformProvider) (client.Status, error) {
	tf, err := Init(providerTerraform)
	if err != nil {
		return client.StatusNotFound, err
	}

	publicKeyBase, err := ssh.GetPublicKeyBase(providerTerraform.Config.MachineFolder)
	if err != nil {
		return client.StatusNotFound, err
	}

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyBase)
	if err != nil {
		return client.StatusNotFound, err
	}
	err = tf.Refresh(context.Background(),
		tfexec.Lock(false),
		tfexec.State(providerTerraform.State),
		tfexec.Var("disk_image="+providerTerraform.Config.DiskImage),
		tfexec.Var("disk_size="+providerTerraform.Config.DiskSizeGB),
		tfexec.Var("instance_type="+providerTerraform.Config.MachineType),
		tfexec.Var("machine_name="+providerTerraform.Config.MachineID),
		tfexec.Var("region="+providerTerraform.Config.Zone),
		tfexec.Var("ssh_key="+string(publicKey)),
	)
	if err != nil {
		return client.StatusNotFound, err
	}

	state, err := tf.ShowStateFile(
		context.Background(),
		providerTerraform.State,
	)
	if err != nil {
		return client.StatusNotFound, err
	}

	if state.Values == nil {
		return client.StatusNotFound, nil
	}
	if state.Values.Outputs != nil {
		return client.StatusRunning, nil
	}

	return client.StatusBusy, nil
}
