package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wgmouton/.dotfiles/manager/commands"
	"github.com/wgmouton/.dotfiles/manager/state"
	"github.com/wgmouton/.dotfiles/manager/types"
	"gopkg.in/yaml.v3"
)

func init() {

	files, _ := filepath.Glob("../tools/*/install.yaml")

	var scripts []types.InstallScriptDefinition

	for _, filePath := range files {
		fileBytes, _ := os.ReadFile(filePath)

		script := types.InstallScriptDefinition{}
		if err := yaml.Unmarshal(fileBytes, &script); err != nil {
			fmt.Println(err)
			continue
		}

		toolPath, _ := filepath.Abs(filepath.Dir(filePath))
		script.ToolPath = toolPath
		scripts = append(scripts, script)
	}

	state.SetToolList(scripts)

	// err := filepath.Walk(".",
	// 	func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		fmt.Println(path, info.Size())
	// 		return nil
	// 	})
	// if err != nil {
	// 	log.Println(err)
	// }

	// t := T{}

	// err := yaml.Unmarshal([]byte(data), &t)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t:\n%v\n\n", t)

	// d, err := yaml.Marshal(&t)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t dump:\n%s\n\n", string(d))

	// m := make(map[interface{}]interface{})

	// err = yaml.Unmarshal([]byte(data), &m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- m:\n%v\n\n", m)

	// d, err = yaml.Marshal(&m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- m dump:\n%s\n\n", string(d))

	// cobra.OnInitialize(initConfig)

	// // Here you will define your flags and configuration settings.
	// // Cobra supports persistent flags, which, if defined here,
	// // will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra-example.yaml)")

	// // Cobra also supports local flags, which will only run
	// // when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
