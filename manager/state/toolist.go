package state

import (
	"context"
	"fmt"

	"github.com/wgmouton/.dotfiles/manager/types"
)

func init() {

	fmt.Println("Building ToolList")

}

var toolList map[int][]types.InstallScriptDefinition = make(map[int][]types.InstallScriptDefinition)

func SetToolList(scripts []types.InstallScriptDefinition) {

	for _, script := range scripts {
		script.InitChannels(context.TODO())
		toolList[script.Stage] = append(toolList[script.Stage], script)
	}

}

func GetToolList() map[int][]types.InstallScriptDefinition {
	return toolList
}
