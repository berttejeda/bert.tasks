package main

import (
	"fmt"
	"github.com/berttejeda/bert.yamlcli/ansible"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	// First positional parameter is the path to the playbook
	if len(os.Args) > 0{
		fmt.Println("ok")
	}
	var playbook string = os.Args[1]
	if strings.HasSuffix(playbook, ".yaml") || strings.HasSuffix(playbook, ".yml") {
  	// Remove the first positional parameter by re-slicing os.Args
  	os.Args = append(os.Args[:1], os.Args[2:]...)
	} else {
		playbook = "Taskfile.yaml"
	}
	cli, cliArguments := ansible.MakeCLIFromAnsiblePlaybook(playbook)
	logger.Debug(cliArguments)
	fmt.Println(cli)
}
