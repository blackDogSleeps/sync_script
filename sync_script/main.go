package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

func getProdAppsFromLocalVariable() []string {
	// Reading from 'vars.sh'
	envFilePath := filepath.Join("apps", "_sh", "vars.sh")
	envFile, envErr := os.Open(envFilePath)

	defer envFile.Close()

	if envErr != nil {
		fmt.Println("Couldn't read file:", envErr)
		os.Exit(0)
	}

	envReader := bufio.NewReader(envFile)
	newLine := uint8(10)
	buffer, err := envReader.ReadBytes(newLine)

	if err != nil {
		fmt.Println("Can't read file: ", err)
		os.Exit(0)
	}

	var prodApps []string
	length := len(buffer)

	for length > 0 {
		buffer, _ = envReader.ReadBytes(newLine)
		length = len(buffer)

		if length < 2 {
			continue
		}

		// If line starts with "
		if buffer[0] == 34 {
			prodApps = append(prodApps, string(buffer[1 : length - 2]))
		}
	}

	return prodApps
}

func copyFilesHelper(fileName string, path string) {
	fmt.Printf("Copying '%s' ", fileName)
	fmt.Printf("to '%s'", path)

	rmErr := os.RemoveAll(path)

	if rmErr != nil {
		fmt.Println("rm error: ", rmErr)
		os.Exit(0)
	}

	cmd := exec.Command("cp", "-rf", fileName, path)
	cmdRes, cmdErr := cmd.Output()

	if cmdErr != nil {
		fmt.Println("Something's not right", cmdErr)
		os.Exit(0)
	}

	fmt.Println(string(cmdRes))
}

func copyFiles(projectName string, filesToCopy []string, directories []string) {
	for i := 0; i < len(filesToCopy); i++ {
		copyFilesHelper(filesToCopy[i], filepath.Join("apps", projectName, filesToCopy[i]))
	}

	for i := 0; i < len(directories); i++ {
		copyFilesHelper(directories[i], filepath.Join("apps", projectName, directories[i]))
	}
}

func buildAssets(arg string) {
	npmCommand := exec.Command("npm", "run", arg)
	npmOut, npmErr := npmCommand.Output()

	fmt.Printf("building %v\n", arg)
	fmt.Println(string(npmOut))

	if npmErr != nil {
		fmt.Println("error: ", npmErr)
		os.Exit(0)
	}
}

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Not enough arguments")
		fmt.Println("Example: -s beta")
		os.Exit(0)
	}

	sOption := args[1]

	if sOption != "-s" {
		fmt.Println("The first argument should be '-s'")
		os.Exit(0)
	}

	projectName := args[2]
	prodApps := getProdAppsFromLocalVariable()
	filesToCopy := []string{
		//"src",
		".browserslistrc",
		".eslintrc.js",
		"cypress.json",
		"jest.config.js",
		"package.json",
		"package-lock.json",
		"postcss.config.js",
		"vue.config.js",
	}

	directories := []string{
		"src",
		"public/pdf",
		"public/tinymce",
	}

	allApps := slices.Concat(prodApps, []string{"beta", "stage"})

	if !slices.Contains(allApps, projectName) && projectName != "prod" {
		fmt.Println("No such project")
		os.Exit(0)
	}

	buildAssets("tailwind-build")
	buildAssets("svg")

	successMessage := "Hooray! It's done!"

	if projectName != "prod" {
		copyFiles(projectName, filesToCopy, directories)
		fmt.Println(successMessage)
		return
	}

	for i := 0; i < len(prodApps); i++ {
		copyFiles(prodApps[i], filesToCopy, directories)
	}

	fmt.Println(successMessage)
}
