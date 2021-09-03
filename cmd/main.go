package main

import (
	"github.com/adeturner/azureBilling/cli"
	"github.com/adeturner/azureBilling/observability"
)

func CliEntry() {
	err := cli.TopLevelCmd.Execute()
	if err != nil && err.Error() != "" {
		observability.Error(err.Error())
	}
}

func main() {
	CliEntry()
}
