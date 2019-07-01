package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
)

var directoriesToSkip = []string{
	"node_modules",
	".git",
}

var startingDirectory string
var searchRecursively bool
var whatToSearch string
var extensionsToSearch string
var whatToSearchFor string
var replacement string

func main() {
	app := cli.NewApp()
	app.Name = "goreplace"
	app.Usage = "CLI application for replacing text in files"

	app.Action = func(c *cli.Context) error {
		promptForStartingDirectory()
		// promptForRecursiveSearch()
		promptForWhatToSearch()

		if whatToSearch == "Specific extensions" {
			promptForExtensionsToSearch()
		}

		promptForWhatToSearchFor()
		promptForReplacement()

		fmt.Printf("\n\nSTARTING: \n")
		perform()

		return nil
	}

	app.Run(os.Args)
}

func promptForStartingDirectory() {
	prompt := promptui.Prompt{
		Label: "Enter the directory to start searching in: ",
	}

	startingDirectory, _ = prompt.Run()
}

func promptForRecursiveSearch() {
	selection := ""
	searchRecursively = false

	prompt := promptui.Select{
		Label: "Search for files recursively? ",
		Items: []string{
			"Yes",
			"No",
		},
	}

	_, selection, _ = prompt.Run()

	if selection == "Yes" {
		searchRecursively = true
	}
}

func promptForWhatToSearch() {
	prompt := promptui.Select{
		Label: "Search: ",
		Items: []string{
			"All files",
			"Specific extensions",
		},
	}

	_, whatToSearch, _ = prompt.Run()
}

func promptForExtensionsToSearch() {
	prompt := promptui.Prompt{
		Label: "Extensions to search (separated by comma, include period): ",
	}

	extensionsToSearch, _ = prompt.Run()
}

func promptForWhatToSearchFor() {
	prompt := promptui.Prompt{
		Label: "Enter what to search for (regex): ",
	}

	whatToSearchFor, _ = prompt.Run()
}

func promptForReplacement() {
	prompt := promptui.Prompt{
		Label: "Enter replacement string: ",
	}

	replacement, _ = prompt.Run()
}

func perform() error {
	extensions := strings.Split(extensionsToSearch, ",")

	err := filepath.Walk(startingDirectory, func(path string, info os.FileInfo, err error) error {
		var fileContents []byte
		continueProcessing := false

		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			for _, dirToSkip := range directoriesToSkip {
				if info.Name() == dirToSkip {
					fmt.Printf("Skipping directory %s\n", info.Name())
					return nil
				}
			}
		}

		if whatToSearch == "Specific extensions" {
			if !info.IsDir() {
				for _, extension := range extensions {
					fileExtension := filepath.Ext(info.Name())

					if fileExtension == extension {
						continueProcessing = true
					}
				}

				if !continueProcessing {
					return nil
				}
			} else {
				return nil
			}
		}

		//fullPath := filepath.Join(path, info.Name())

		if fileContents, err = ioutil.ReadFile(path); err != nil {
			fmt.Printf("Error reading file %q: %v\n", path, err)
			return err
		}

		r := regexp.MustCompile(whatToSearchFor)

		if r.Match(fileContents) {
			newFileContents := r.ReplaceAll(fileContents, []byte(replacement))

			if err = ioutil.WriteFile(path, newFileContents, info.Mode()); err != nil {
				fmt.Printf("Error writing updated file %s: %v\n", path, err)
				return err
			}

			fmt.Printf("File %s updated\n", path)
		}

		return nil
	})

	return err
}
