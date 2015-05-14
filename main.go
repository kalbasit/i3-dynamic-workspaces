package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kalbasit/i3-dynamic-workspaces/dmenu"
	"github.com/proxypoke/i3ipc"
)

var (
	i3  *i3ipc.IPCSocket
	err error

	flagSwitchworkspace          bool
	flagMoveContainerToWorkspace bool
)

func init() {
	flag.BoolVar(&flagSwitchworkspace,
		"switch-workspace",
		false,
		"Switch to an existing workspace")

	flag.BoolVar(&flagMoveContainerToWorkspace,
		"move-container-to-workspace",
		false,
		"Move focused container to workspace")

	flag.Usage = usage
}

func main() {
	if i3, err = i3ipc.GetIPCSocket(); err != nil {
		panic(err)
	}

	flag.Parse()

	if (!flagSwitchworkspace && !flagMoveContainerToWorkspace) ||
		(flagSwitchworkspace && flagMoveContainerToWorkspace) {
		usage()
	}

	// Let's get the workspace to move to with dmenu
	workspace := getWorkspace()

	if flagMoveContainerToWorkspace {
		moveContainerToWorkspace(workspace)
	} else if flagSwitchworkspace {
		selectWorkspace(workspace)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s [flags]\n", path.Base(os.Args[0]))
	flag.VisitAll(func(fg *flag.Flag) {
		fmt.Printf("--%s=\"%s\"\n\t%s\n", fg.Name, fg.DefValue,
			strings.Replace(fg.Usage, "\n", "\n\t", -1))
	})
	os.Exit(1)
}

func moveContainerToWorkspace(workspace string) {
	command := fmt.Sprintf("move container to workspace %s", workspace)
	if _, err := i3.Command(command); err != nil {
		panic(err)
	}
}

func getWorkspace() string {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		panic(err)
	}

	items := make([]string, 0, len(workspaces))

	for _, workspace := range workspaces {
		item := workspace.Name

		items = append(items, item)
	}

	workspace := dmenu.Run(items)

	return workspace
}

func selectWorkspace(workspace string) {
	command := fmt.Sprintf("workspace %s", workspace)
	if _, err := i3.Command(command); err != nil {
		panic(err)
	}
}
