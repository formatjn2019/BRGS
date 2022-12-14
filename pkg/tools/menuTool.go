package tools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// 命令行菜单
func CommandMenu(canBack bool, menuItems ...string) int {
	if canBack {
		menuItems = append(menuItems, "返回上级")
	}
	menuSelect := func(menuItems ...string) int {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			for index, menuItem := range menuItems {
				fmt.Printf("%5d.%80s\n", index+1, menuItem)
			}
			scanner.Scan()
			text := scanner.Text()
			if index, err := strconv.Atoi(text); err == nil {
				if index <= len(menuItems) && index > 0 {
					return index - 1
				} else {
					fmt.Println("输入无效，请重新输入")
				}
			}
		}
	}
	//索引转换
	index := menuSelect(menuItems...)
	if canBack {
		return (index+1)%len(menuItems) - 1
	}
	return index
}
