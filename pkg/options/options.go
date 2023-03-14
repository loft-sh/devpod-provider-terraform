package options

import (
	"fmt"
	"os"
)

const (
	REGION            = "REGION"
	INSTANCE_TYPE     = "INSTANCE_TYPE"
	IMAGE_DISK        = "IMAGE_DISK"
	DISK_SIZE         = "DISK_SIZE"
	TERRAFORM_PROJECT = "TERRAFORM_PROJECT"
)

type Options struct {
	DiskImage     string
	DiskSizeGB    string
	MachineFolder string
	MachineID     string
	MachineType   string
	Zone          string
}

func ConfigFromEnv() (Options, error) {
	return Options{
		MachineType: os.Getenv(INSTANCE_TYPE),
		DiskImage:   os.Getenv(IMAGE_DISK),
		DiskSizeGB:  os.Getenv(DISK_SIZE),
		Zone:        os.Getenv(REGION),
	}, nil
}

func FromEnv() (*Options, error) {
	retOptions := &Options{}

	var err error

	retOptions.MachineID, err = FromEnvOrError("MACHINE_ID")
	if err != nil {
		return nil, err
	}
	// prefix with devpod-
	retOptions.MachineID = "devpod-" + retOptions.MachineID

	retOptions.MachineFolder, err = FromEnvOrError("MACHINE_FOLDER")
	if err != nil {
		return nil, err
	}

	retOptions.MachineType, err = FromEnvOrError("INSTANCE_TYPE")
	if err != nil {
		return nil, err
	}

	retOptions.DiskImage, err = FromEnvOrError("IMAGE_DISK")
	if err != nil {
		return nil, err
	}

	retOptions.DiskSizeGB, err = FromEnvOrError("DISK_SIZE")
	if err != nil {
		return nil, err
	}

	retOptions.Zone, err = FromEnvOrError("REGION")
	if err != nil {
		return nil, err
	}

	return retOptions, nil
}

func FromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf(
			"couldn't find option %s in environment, please make sure %s is defined",
			name,
			name,
		)
	}

	return val, nil
}
