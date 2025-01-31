package main

import (
	"bufio"
	"fmt"
	"github.com/berttejeda/bert.tasks/lib"
	"github.com/berttejeda/bert.yamlcli/ansible"
	logger "github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	// Set log flags to include date, time, and short file name
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	x := "EX!"

	// Use a defer function to recover from panics and log the error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Something went wrong: '%s', error was %s", x, r)
			os.Exit(1)
		}
	}()

	var playbook string
	var defaultPlaybook string = "Taskfile.yaml"
	// First positional parameter is the path to the playbook
	if len(os.Args) > 1 {
		playbook = os.Args[1]
		if strings.HasSuffix(playbook, ".yaml") || strings.HasSuffix(playbook, ".yml") {
			// Remove the first positional parameter by re-slicing os.Args
			if strings.HasPrefix(playbook, "~/") {
				dirname, _ := os.UserHomeDir()
				playbook = filepath.Join(dirname, playbook[2:])
			}
			os.Args = append(os.Args[:1], os.Args[2:]...)
		} else {
			playbook = defaultPlaybook
		}
	} else {
		playbook = defaultPlaybook
	}

	cmd, cmdOptions, ansibleCLI, ansibleCLIOptions, ansibleScriptWrapperFile, echoOn := ansible.MakeCLIFromAnsiblePlaybook(playbook, os.Args)
	logger.Debug(cmd, cmdOptions, ansibleCLI, ansibleCLIOptions)

	if echoOn {
		fmt.Printf(ansibleCLI)
	} else {

		ansibleScriptFile, err := lib.CreateFile(ansibleScriptWrapperFile, ansibleCLI)
		if err != nil {
			logger.Fatal(fmt.Sprintf("error: %w", err))
		}
		ansibleCLIInstance := exec.Command("bash", ansibleScriptFile)

		// Create pipes for stdout and stderr
		stdout, err := ansibleCLIInstance.StdoutPipe()
		if err != nil {
			logger.Error(fmt.Sprintf("error creating stdout pipe: %v\n", err))
		}

		stderr, err := ansibleCLIInstance.StderrPipe()
		if err != nil {
			logger.Error(fmt.Sprintf("error creating stderr pipe: %v\n", err))
		}

		// Start the command
		if err := ansibleCLIInstance.Start(); err != nil {
			logger.Fatal(fmt.Sprintf("Error starting command: %v\n", err))
		}

		// Stream stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Printf("%s\n", scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				logger.Debug(fmt.Sprintf("Error reading stdout: %v\n", err))
			}
		}()

		// Stream stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Printf("%s\n", scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				logger.Debug(fmt.Sprintf("Error reading stderr: %v\n", err))
			}
		}()

		// Wait for the command to complete
		if err := ansibleCLIInstance.Wait(); err != nil {
			logger.Warning(fmt.Sprintf("Error waiting for command to finish: %v\n", err))
		} else {
			fmt.Println("Command executed successfully.")
		}

	}

}
