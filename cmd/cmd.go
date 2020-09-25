package cmd

import (
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use: "rman",
}

func Execute() {
	_ = rootCommand.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	toml.NewEncoder(os.Stdout)
}

func initConfig() {

}
