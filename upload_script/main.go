package main

import (
	"bufio"
	"fmt"
	"os"
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

	var errors []string;
	if len(userName) < 1 {
		errors = append(errors, "Не указан SERVER_USER_NAME в .env")
		//fmt.Println("Не указан SERVER_USER_NAME в .env")
		//os.Exit(0)
	}

	if len(pubKeyPath) < 1 {
		errors = append(errors, "Не указан SERVER_PUB_KEY_PATH в .env")
		//fmt.Println("Не указан SERVER_PUB_KEY_PATH в .env")
		//os.Exit(0)
	}

	if len(appName) < 1 {
		errors = append(errors, "Не указан SERVER_APP_NAME в .env")
		//fmt.Println("Не указан SERVER_APP_NAME в .env")
		//os.Exit(0)
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

func main() {
	envVars := getLocalEnvVariables()
	
	fmt.Println(envVars);

	 
	//if [[ ! $SERVER_USER_NAME ]]
	//then
	//    echo "Не указан SERVER_USER_NAME в .env "
	//    exit
	//fi
	//if [[ ! $SERVER_PUB_KEY_PATH ]]
	//then
	//    echo "Не указан SERVER_PUB_KEY_PATH в .env "
	//    exit
	//fi
	//if [[ ! $SERVER_APP_NAME ]]
	//then
	//    echo "Не указан SERVER_APP_NAME в .env "
	//    exit
	//fi

}
