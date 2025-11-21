package main

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
			fmt.Println("Can't read your input: ", scanErr)
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


func makeDistArchive() {
	fmt.Println("Архивируем dist папку")
	dirErr := os.Chdir(filepath.Join("dist"))
	
	if dirErr != nil {
		fmt.Println("Some trouble with os.Chdir: ", dirErr)
		os.Exit(0)
	}

	// Remove old 'dist.zip'
	rmCmdErr := os.Remove("dist.zip")

	if rmCmdErr != nil {
		fmt.Println("Couldn't delete 'dist.zip': ", rmCmdErr)
		os.Exit(0)
	}

	// Make new archive
	makeNewArchive()
}

func makeNewArchive() {
	zipPath := filepath.Join("..", "dist.zip")
	zipFile, err := os.Create(zipPath)

	if err != nil {
		fmt.Println("Couldn't create 'dist.zip': ", err)
		os.Exit(0)
	}
	
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	writeErr := writer.AddFS(os.DirFS("."))

	if writeErr != nil {
		fmt.Println("Couldn't zip the files: ", writeErr)
		os.Exit(0)
	}

	currDir, _ := os.Getwd()
	newZipPath := filepath.Join(currDir, "dist.zip")
	cpErr := os.Rename(zipPath, newZipPath)

	if cpErr != nil {
		fmt.Println("cpErr: ", cpErr)
		os.Exit(0)
	}
}

func uploadArchiveToServer(userName string, pubKeyPath string, appPath string) {
	// Command we're translating:
	//scp -i $SERVER_PUB_KEY_PATH -P 58533 dist.zip $SERVER_USER_NAME@88.99.160.4:$_appPath
	serverLogin := fmt.Sprintf("%s@88.99.160.4:%s", userName, appPath)
	ctx := context.Background()
	uploadCmd := exec.CommandContext(
		ctx, 
		"scp", 
		"-i", 
		pubKeyPath,
		"-P", 
		"58533",
		"dist.zip",
		serverLogin,
	)

	uploadCmd.Stdin = os.Stdin
	uploadCmd.Stdout = os.Stdout
	uploadCmd.Stderr = os.Stderr
	uploadErr := uploadCmd.Run()

	if uploadErr != nil {
		fmt.Println("scp command didn't fly: ", uploadErr)
		os.Exit(0)
	}
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

	// The ssh command we're translating:
	//ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'rm -rf '$_appPath'/* '$_appPath'/.htaccess'
	fmt.Println("Удаляем файлы")
	serverCommands = fmt.Sprintf("rm -rf %s/* %s/.htaccess", appPath, appPath)
	execCommand(envVars["userName"], serverCommands)
	makeDistArchive()
	fmt.Println("Теперь копируем архив на сервер")
	uploadArchiveToServer(envVars["userName"], envVars["pubKeyPath"], appPath)

	// The ssh command we're translating:
	//ssh $SERVER_USER_NAME@88.99.160.4 -p 58533 'bash -c "cd '$_appPath'/ && unzip dist.zip && rm -f dist.zip"'
	fmt.Println("Распакуем его и удаляем архив")
	serverCommands = fmt.Sprintf("bash -c 'cd %s/ && unzip dist.zip && rm -f dist.zip && chmod -R 775 *'", appPath)
	execCommand(envVars["userName"], serverCommands)
	fmt.Println("Идём и проверяем что всё открывается!")
	fmt.Printf("https://%s.app.motivationportal.ru\n", envVars["appName"])
	dirErr := os.Chdir(filepath.Join(".."))
	
	if dirErr != nil {
		fmt.Println("Some trouble with os.Chdir: ", dirErr)
		os.Exit(0)
	}
}
