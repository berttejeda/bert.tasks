# Overview

*bert.tasks* is an ansible wrapper that uses an ansible playbook 
file to define its command-line parameters.

If no playbook file is specified, the app will search for 'Taskfile.yaml' in the current working directory.

Consult the data structure of the included ansible playbook: [Taskfile.yaml](Taskfile.yaml).

# Usage Example

`go run main.go run --help`

