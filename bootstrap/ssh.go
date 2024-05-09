package bootstrap

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"mutilcmd/global"
	"os"
)

func InitializeSSH() (*ssh.Client, error) {
	// 尝试使用私钥进行身份验证	
	privateKeyBytes, err := os.ReadFile(global.App.Config.App.PrivateKey)
	if err == nil {
		// global.App.Log.Info("try to use private key for authentication")
		privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err == nil {
			config := &ssh.ClientConfig{
				User: global.App.Config.App.User,
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(privateKey),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			client, err := ssh.Dial("tcp", global.App.Config.App.Host+":22", config)
			if err == nil {
				// global.App.Log.Info(fmt.Sprintf("success to use private key for authentication，host：%s", global.App.Config.App.Host))
				return client, nil
			} else {
				global.App.Log.Error(fmt.Sprintf("cannot connect to host：%s, host unreachable！", global.App.Config.App.Host))
				return nil, err
			}
		}
	} else {
		global.App.Log.Error(fmt.Sprintf("cannot read private key：%s", global.App.Config.App.PrivateKey))
		return nil, err
	}

	// 使用密码进行身份验证
	global.App.Log.Info("please input password for authentication")

	config := &ssh.ClientConfig{
		User: global.App.Config.App.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(global.App.Config.App.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", global.App.Config.App.Host+":22", config)
	if err != nil {
		global.App.Log.Error(fmt.Sprintf("cannot connect to host：%s, host unreachable！", global.App.Config.App.Host))
		return nil, err
	}
	// global.App.Log.Info(fmt.Sprintf("success to use password for authentication，host：%s", global.App.Config.App.Host))
	return client, nil
}
