# TERRAFORM Provider for DevPod

[![Join us on Slack!](docs/static/media/slack.svg)](https://slack.loft.sh/) [![Open in DevPod!](https://devpod.sh/assets/open-in-devpod.svg)](https://devpod.sh/open#https://github.com/loft-sh/devpod-provider-terraform)

## Getting started

The provider is available for auto-installation using 

```sh
devpod provider add terraform
devpod provider use terraform
```

Follow the on-screen instructions to complete the setup.

Needed variables will be:

- TERRAFORM_PROJECT
- REGION

TERRAFORM_PROJECT points to a git repo or directory where the terraform project
that defines the infra is stored.

In this repo it would point to: `./examples/terraform-aws/`

### Creating your first devpod env with terraform

After the initial setup, just use:

```sh
devpod up .
```

You'll need to wait for the machine and environment setup.

### Notes

With the terraform provider, all the power is in the terraform project. So it
will be there where you will place your defaults for

- DISK_SIZE
- IMAGE_DISK
- INSTANCE_TYPE

Keep also in mind that **stop/start is not supported right now on terraform provider**
So the right thing to do is to handle data saving inside your terraform code
(eg. use external data buckets)

### Customize the VM Instance

This provides has the seguent options

|    NAME           | REQUIRED |          DESCRIPTION                  |         DEFAULT         |
|-------------------|----------|---------------------------------------|-------------------------|
| DISK_SIZE                 | false    | The disk size to use.                 | 40  |
| IMAGE_DISK                | false    | The disk image to use.                |     |
| INSTANCE_TYPE             | false    | The machine type to use.              |     |
| REGION                    | true     | The cloud region to create the        |     |
|                           |          | VM in. E.g. us-west-1                 |     |
| TERRAFORM_PROJECT         | true     | The path or repo where the            |     |
|                           |          | terraform files are. E.g.             |     |
|                           |          | ./examples/terraform or               |     |
|                           |          | https://github.com/examples/terraform |     |

Options can either be set in `env` or using for example:

```sh
devpod provider set-options -o IMAGE_DISK=my-custom-ami
devpod provider set-options -o INSTANCE_TYPE=t2.micro
devpod provider set-options -o REGION=us-west-2
```
