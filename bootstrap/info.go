package bootstrap

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
	"golang.org/x/term"
	"syscall"
)

type ReadInfo struct {
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	PrivateKey string `yaml:"private_key"`
	Password   string `yaml:"password"`
}

func InitializeInfo() {

	config := ReadInfo{}
	config.readConfig()

	err := writeConfigToYAML(config)
	if err != nil {
		fmt.Println("write config to yaml failed:", err)
		return
	}

	fmt.Println("successfully written to config.yaml file.")
}

func (config *ReadInfo) readConfig() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Host address: ")
	host, _ := reader.ReadString('\n')
	config.Host = strings.TrimSpace(host)

	fmt.Print("Username: ")
	user, _ := reader.ReadString('\n')
	config.User = strings.TrimSpace(user)

	fmt.Print("private key path: ")
	privateKey, _ := reader.ReadString('\n')
	config.PrivateKey = strings.TrimSpace(privateKey)

	fmt.Print("Password: ")
	passwordBytes, _ := term.ReadPassword(int(syscall.Stdin))
	config.Password = string(passwordBytes)
}

func writeConfigToYAML(config ReadInfo) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshDir := homeDir + "/.mutilcmd/config"
	err = os.MkdirAll(sshDir, 0700)
	if err != nil {
		return err
	}

	filePath := sshDir + "/config.yaml"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlData, err := yaml.Marshal(map[string]interface{}{
		"app": config,
		"log": map[string]interface{}{
			"level":       "info",
			"root_dir":    "./runtime/logs",
			"filename":    "app.log",
			"format":      "json",
			"show_line":   true,
			"max_backups": 3,
			"max_size":    500,
			"max_age":     28,
			"compress":    true,
		},
	})
	if err != nil {
		return err
	}

	// yamlDataString := string(yamlData)
	// yamlDataString = strings.Replace(yamlDataString, fmt.Sprintf("private_key: %s", config.PrivateKey), fmt.Sprintf("private_key: '%s'", config.PrivateKey), 1)

	// _, err = file.WriteString(yamlDataString)
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	return nil
}