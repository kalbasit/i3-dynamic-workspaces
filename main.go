package main

import (
	"flag"

	"fmt"

	"github.com/eMxyzptlk/i3-dynamic-workspaces/dmenu"
	"github.com/proxypoke/i3ipc"

	"log"

	"os"

	"path"

	"strings"
)

var (
	i3  *i3ipc.IPCSocket
	err error

	flagAddworkspace    bool
	flagSwitchworkspace bool
)

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage: %s [flags]\n", path.Base(os.Args[0]))
	flag.VisitAll(func(fg *flag.Flag) {
		fmt.Printf("--%s=\"%s\"\n\t%s\n", fg.Name, fg.DefValue,
			strings.Replace(fg.Usage, "\n", "\n\t", -1))
	})
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flagAddworkspace, "add-workspace", false, "Add a new workspace")
	flag.BoolVar(&flagSwitchworkspace, "switch-workspace", false, "Switch to an existing workspace")

	flag.Usage = usage
}

func addWorkspace() {
	workspace := dmenu.Run([]string{})
	selectWorkspace(workspace)
}

func switchWorkspace() {
	workspaces, err := i3.GetWorkspaces()

	if err != nil {
		log.Fatalln(err)
	}

	items := make([]string, 0, len(workspaces))

	for _, workspace := range workspaces {
		item := workspace.Name

		items = append(items, item)
	}

	workspace := dmenu.Run(items)
	selectWorkspace(workspace)

}

func selectWorkspace(workspace string) {
	i3.Command(fmt.Sprintf("workspace %s", workspace))
}

func main() {
	i3, err = i3ipc.GetIPCSocket()
	if err != nil {
		log.Fatalln(err)
	}

	flag.Parse()

	if flagAddworkspace {
		addWorkspace()
	} else if flagSwitchworkspace {
		switchWorkspace()
	}
}
