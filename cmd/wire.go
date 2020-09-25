package cmd

import "github.com/spf13/cobra"

func init() {
	rootCommand.AddCommand(wireCommand)
}

var wireCommand = &cobra.Command{
	Use:   "wire [path]",
	Short: "Adds replace directives to go.mod for sibling projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		workDir := "."
		if len(args) == 1 {
			workDir = args[0]
		}

		return wireSiblingsInPath(workDir)
	},
}
