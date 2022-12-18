package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wgmouton/.dotfiles/manager/state"
)

var listToolsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tools",
	Long:  `This subcommand says hello`,
	Run: func(cmd *cobra.Command, args []string) {
		state := state.GetToolList()
		fmt.Printf("%+v", state)
	},
}

var inspectToolsCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect specific tool",
	Long:  `This subcommand says hello`,
	Run: func(cmd *cobra.Command, args []string) {
		state := state.GetToolList()
		fmt.Printf("%+v", state)
	},
}

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "tools",
	Long:  `This subcommand says hello`,
}

func init() {
	toolsCmd.AddCommand(listToolsCmd, inspectToolsCmd)
	RootCmd.AddCommand(toolsCmd)
}
