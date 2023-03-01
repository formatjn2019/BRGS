package main

import (
	"BRGS/management/menu"
)

func main() {
	// commands := management.METHODS_ALL
	// commandArr := make([]string, 0)
	// for _, command := range commands {
	// 	commandArr = append(commandArr, command.String())
	// }
	// management.ShowCommands()
	// for {
	// 	index := management.CommandMenu(commandArr...)
	// 	if index >= 0 {
	// 		commands[index].Execute()
	// 	} else {
	// 		break
	// 	}
	// }
	menu.StartMenu()

}
