package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/berttejeda/bert.yamlcli/ansible"
	logger "github.com/sirupsen/logrus"
	"os"
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
	app := kingpin.New("", "")
	debug := app.Flag("debug", "Enable Debug Logging").Bool()
	cli, cliArguments := ansible.MakeCLIFromAnsiblePlaybook(app, playbook)
	if *debug {
		logger.SetLevel(logger.DebugLevel)
	} else {
		logger.SetLevel(logger.InfoLevel)
	}
	logger.Debug(*debug)
	logger.Debug(*cliArguments["foo"].(*string))
	logger.Debug(*cliArguments["bar"].(*string))
	logger.Debug(cli)
}
