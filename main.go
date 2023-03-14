package main

import (
	"github.com/loft-sh/devpod-provider-terraform/cmd"
)

func main() {
	cmd.Execute()
}

/*
TODO: command == ssh to output IP of create
TODO: stop == scale to 0 the instance
TODO: start == scale back to 1 the instance
TODO: status == ??? we check the plan and parse it from json?
TODO: Terraform was run at least once successful and not scaled to 0 => Running
TODO: Terraform was run at least once successful and scaled to 0 => Stopped
TODO: Terraform was not run successful => NotFound
*/
