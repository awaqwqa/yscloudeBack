package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// CmdController 结构体用于管理和处理命令
type CmdController struct {
	commands sync.Map
}

func NewCmdController() *CmdController {
	return &CmdController{}
}

// RegisterCommand 方法用于注册命令处理函数
func (c *CmdController) RegisterCommand(name string, handler func(args []string)) {
	c.commands.Store(name, handler)
}

// ExecuteCommand 方法用于执行命令
func (c *CmdController) ExecuteCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	cmdName := parts[0]
	args := parts[1:]

	if cmd, ok := c.commands.Load(cmdName); ok {
		cmd.(func(args []string))(args)
	} else {
		fmt.Println("Unknown command:", cmdName)
	}
}
func DbCmdHandle(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "User":
		if len(args) < 2 {

		}
	}
	fmt.Println()
}
func (c *CmdController) Init() {
	c.RegisterCommand("DB", DbCmdHandle)
}

// Listen 方法用于持续监听终端输入
func (c *CmdController) Listen() {
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			fmt.Print("> ")
			if scanner.Scan() {
				input := scanner.Text()
				c.ExecuteCommand(input)
			}
		}

	}()

}
