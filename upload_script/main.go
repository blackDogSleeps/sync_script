package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
)

func getLocalEnvVariables() map[string]string {
	envFilePath := ".env.local"
	envFile, envErr := os.Open(envFilePath)

	defer envFile.Close()

	if envErr != nil {
		fmt.Println("Couldn't read file")
		os.Exit(0)
	}

	envReader := bufio.NewReader(envFile)
	newLine := uint8(10)
	buffer, err := envReader.ReadBytes(newLine)

	if err != nil {
		fmt.Println("Can't read file: ", err)
		os.Exit(0)
	}

	// Variables we expect:
	//SERVER_USER_NAME={userName}
	//SERVER_PUB_KEY_PATH={pubKeyPath}
	//SERVER_APP_NAME={appName}
	var userName, pubKeyPath, appName string
	length := len(buffer)

	for length > 0 {
		buffer, _ = envReader.ReadBytes(newLine)
		length = len(buffer)

		if length < 2 {
			continue
		}

		if string(buffer)[7:11] == "USER" {
			userName = string(buffer[17 : length-1])
		}
		if string(buffer)[7:10] == "PUB" {
			pubKeyPath = string(buffer[20 : length-1])
		}

		if string(buffer)[7:10] == "APP" {
			appName = string(buffer[16 : length-1])
		}
	}

	var errors []string

	if len(userName) < 1 {
		errors = append(errors, "Не указан SERVER_USER_NAME в .env")
	}

	if len(pubKeyPath) < 1 {
		errors = append(errors, "Не указан SERVER_PUB_KEY_PATH в .env")
	}

	if len(appName) < 1 {
		errors = append(errors, "Не указан SERVER_APP_NAME в .env")
	}

	if len(errors) > 0 {
		for message := range slices.Values(errors) {
			fmt.Println(message)
		}
		os.Exit(0)
	}

	return map[string]string{
		"userName":   userName,
		"pubKeyPath": pubKeyPath,
		"appName":    appName,
	}
}

func execCommand(userName string, serverCommands string) {
	serverLogin := fmt.Sprintf("%s@88.99.160.4", userName)
	fmt.Println(serverLogin)
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "ssh", serverLogin, "-p", "58533", serverCommands)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	sshErr := cmd.Run()
	if sshErr != nil {
		fmt.Println("ssh command didn't fly: ", sshErr)
		os.Exit(0)
	}
}

func getUserInput() {
	var reply string

	fmt.Println("Удалить их и загрузить новые (y/n)?")

	for reply != "y" {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		scanErr := scanner.Err()

		if scanErr != nil {
			fmt.Println("Something wrong with scanner: ", scanErr)
			os.Exit(0)
		}

		reply = scanner.Text()

		if reply == "n" {
			fmt.Println("Возвращайтесь когда будете готовы")
			os.Exit(0)
		}

		if reply != "y" && reply != "n" {
			fmt.Println("Можно выбрать только 'y' или 'n' =)")
		}
	}

	fmt.Println("Аминь!")
}

func main() {
	envVars := getLocalEnvVariables()
	appPath := fmt.Sprintf("/var/www/projects/app/%s", envVars["appName"])
	serverCommands := fmt.Sprintf("cd %s/ && pwd && ls", appPath)

	fmt.Println("Параметры из .env.local")
	fmt.Printf("SERVER_USER_NAME=%s\n", envVars["userName"])
	fmt.Printf("SERVER_PUB_KEY_PATH=%s\n", envVars["pubKeyPath"])
	fmt.Printf("SERVER_APP_NAM=%s\n", envVars["appName"])
	fmt.Printf("appPath: %s\n", appPath)
	fmt.Println("Файлы на сервере:")
	fmt.Println("=======")

	// The ssh command we're translating:
	//ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'cd '$_appPath'/ && pwd && ls'
	execCommand(envVars["userName"], serverCommands)
	getUserInput()

	fmt.Println("Удаляем файлы")
	//ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'rm -rf '$_appPath'/* '$_appPath'/.htaccess'
	serverCommands = fmt.Sprintf("rm -rf %s/* %s/.htaccess", appPath)
	//execCommand(envVars["userName"], serverCommands)

}
