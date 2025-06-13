package main

import (
	"fmt"
	"github.com/berttejeda/bert.tasks/lib"
	"github.com/berttejeda/bert.yamlcli/ansible"
	logger "github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
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
			debug.PrintStack()
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

	cmd, cmdOptions, ansibleCLI, ansibleCLIOptions, ansibleScriptWrapperFile := ansible.MakeCLIFromAnsiblePlaybook(playbook, os.Args)
	logger.Debug(cmd, cmdOptions, ansibleCLI, ansibleCLIOptions)
	var _, echoOn = cmdOptions["--dry-run"]
	if echoOn {
		fmt.Printf(ansibleCLI)
	} else {

		ansibleScriptFile, err := lib.CreateFile(ansibleScriptWrapperFile, ansibleCLI)
		if err != nil {
			logger.Fatal(fmt.Sprintf("error: %w", err))
		}

		ansibleCLIInstance := exec.Command("bash", ansibleScriptFile)

		// The following ensures we are running the command interactively
		ansibleCLIInstance.Stdin = os.Stdin
		ansibleCLIInstance.Stdout = os.Stdout
		ansibleCLIInstance.Stderr = os.Stderr

		// Run the command
		err = ansibleCLIInstance.Run()
		if err != nil {
			fmt.Println("Error running ansible script:", err)
		}
	}

}
