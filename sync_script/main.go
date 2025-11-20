package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
)

func copyFilesHelper(fileName string, path string) {
	fmt.Printf("Copying '%s' ", fileName)
	fmt.Printf("to '%s'", path)

	rmCmd := exec.Command("rm", "-rf", path)
	_, rmErr := rmCmd.Output()

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
	// TODO: import from vars.sh instead
	prodApps := []string{
		"base",
		"vniizht",
		"rolf",
		"azconnect",
		"storeez12",
		"ecco",
		"borjomi",
		"avtosushi",
		"asg",
		"ekonika",
		//"leagueofcare" # клиент ушёл
		// "dreamer" # клиент ушёл
	}

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
