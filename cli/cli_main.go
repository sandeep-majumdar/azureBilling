// https://medium.com/@skdomino/writing-better-clis-one-snake-at-a-time-d22e50e60056
// https://github.com/nanopack/hoarder/blob/master/commands/commands.go
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	body     io.ReadWriter        // what to read/write requests body
	verbose  bool                 // whether to display request info
	insecure bool          = true // whether to ignore cert or not
	showVers bool                 // whether to print version info or not

	// to be populated by linker
	version = "1.0.0"

	// Cobra parameters that are common across collections
	url, id, output, data, file, login, configFile, region, subscription string
	seq, limit, port                                                     int
	y                                                                    bool

	// choreoCmd ...
	TopLevelCmd = &cobra.Command{
		Use:     "azutils",
		Short:   "azutils",
		Long:    ``,
		Example: "",
		RunE:    preFlight,
	}
)

func preFlight(ccmd *cobra.Command, args []string) error {
	// if --version is passed print the version info
	if showVers {
		fmt.Printf("azutils %s\n", version)
	}
	return nil
}

func init() {
	TopLevelCmd.Flags().BoolVarP(&showVers, "version", "v", false, "Display the current version of this CLI")
	TopLevelCmd.AddCommand(azureCmd)
}
