package commands

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "cobra-example",
	Short: "An example of cobra",
	Long: `This application shows how to create modern CLI 
applications in go using Cobra CLI library`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}
