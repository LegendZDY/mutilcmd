package main

import (
	"fmt"
	"flag"
	"os"

	"strings"
	"mutilcmd/global"
	"mutilcmd/bootstrap"
)

func main() {

	// 输入参数
	homeDir := os.Getenv("HOME")
	configPath := homeDir + "/.mutilcmd/config/config.yaml"

	// 解析命令行参数
	flag.Usage = func() {
		fmt.Println("Usage:  mutilcmd is a tool to execute multiple commands on multiple servers.")
		fmt.Println("	When you run mutilcmd without any arguments, it will start a interactive mode to create a config file.")
		fmt.Println("	You can also use the following syntax to execute a command on a server: mutilcmd 'any command'")
		fmt.Println("	For example: mutilcmd 'ls -l'")
		flag.PrintDefaults()
	}
	flag.Parse()
	otherArgs := flag.Args()

	if _, err := os.Stat(configPath); err == nil && len(otherArgs) != 0 {
		// 初始化配置
		global.App.Config.InitializeConfig()
		// 初始化日志
		global.App.Log = bootstrap.InitializeLog()
		global.App.Log.Info("Starting SSH Command")

		client, err := bootstrap.InitializeSSH()
		if err != nil {
			panic(err)
		}
		defer client.Close()

		// 创建新的会话
		session, err := client.NewSession()
		if err != nil {
			global.App.Log.Info(fmt.Sprintf("Failed to create session: %s", err))
		}
		defer session.Close()

		// 执行命令
		command := strings.Join(otherArgs, " ")
		global.App.Log.Info(fmt.Sprintf("Executing command: %s", command))
		output, err := session.CombinedOutput(command)
		if err != nil {
			global.App.Log.Info(fmt.Sprintf("Failed to execute command: %s", err))
		}

		fmt.Println(string(output))
		global.App.Log.Info("SSH Command successfully executed")
	} else {
		bootstrap.InitializeInfo()
	}
}
