package main

import (
	"bufio"
	"fmt"
	"github.com/berttejeda/bert.tasks/lib"
	"github.com/berttejeda/bert.yamlcli/ansible"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// First positional parameter is the path to the playbook
	var playbook string = os.Args[1]
	if strings.HasSuffix(playbook, ".yaml") || strings.HasSuffix(playbook, ".yml") {
		// Remove the first positional parameter by re-slicing os.Args
		if strings.HasPrefix(playbook, "~/") {
			dirname, _ := os.UserHomeDir()
			playbook = filepath.Join(dirname, playbook[2:])
		}
		os.Args = append(os.Args[:1], os.Args[2:]...)
	} else {
		playbook = "Taskfile.yaml"
	}
	// Create the application
	cmd, cmdOptions, ansibleCLI, ansibleCLIOptions, ansibleScriptWrapperFile := ansible.MakeCLIFromAnsiblePlaybook(playbook, os.Args)
	logger.Debug(cmd, cmdOptions, ansibleCLI, ansibleCLIOptions)

	ansibleScriptFile, err := lib.CreateFile(ansibleScriptWrapperFile, ansibleCLI)
	if err != nil {
		fmt.Errorf("Error: %w", err)
	}
	ansibleCLIInstance := exec.Command("bash", ansibleScriptFile)

	// Create pipes for stdout and stderr
	stdout, err := ansibleCLIInstance.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return
	}

	stderr, err := ansibleCLIInstance.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return
	}

	// Start the command
	if err := ansibleCLIInstance.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stdout: %v\n", err)
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading stderr: %v\n", err)
		}
	}()

	// Wait for the command to complete
	if err := ansibleCLIInstance.Wait(); err != nil {
		fmt.Printf("Error waiting for command to finish: %v\n", err)
	} else {
		fmt.Println("Command executed successfully.")
	}

}
