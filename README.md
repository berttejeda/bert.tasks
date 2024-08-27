# Overview

This is a Golang rewrite of the python equivalent [ansible-taskrunner](https://github.com/berttejeda/ansible-taskrunner).

*bert.tasks* is an ansible wrapper that uses an ansible playbook file to define its command-line parameters.

The command-line parameters are used to build a subprocess call to `ansible-playbook`.

The goal is to be able to easily define command-line parameters that translate to ansible variables.

If no playbook file is specified, the app will search for 'Taskfile.yaml' in the current working directory.

Consult the data structure of the included ansible playbook: [Taskfile.yaml](Taskfile.yaml).

# Usage Example

* Show the help output<br />
```bash
go run main.go --help
```
* Show the help output for the `run` subcommand, defined in the default playbook specification [Taskfile.yaml](Taskfile.yaml)<br />
```bash
go run main.go run --help
```
* Show the help output for the `run` subcommand, explicitly specifying the playbook to use as input<br />
```bash
go run main.go Taskfile.yaml run --help
````

