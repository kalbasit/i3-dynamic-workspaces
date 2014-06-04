package main

import (
	"flag"

	"fmt"

	"github.com/proxypoke/i3ipc"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/BurntSushi/wingo/prompt"

	"log"

	"os"

	"path"

	"strings"
)

var (
	x   *xgbutil.XUtil
	i3  *i3ipc.IPCSocket
	err error

	flagAddworkspace bool
)

// response is the callback that gets executed whenever the user hits
// enter (the "confirm" key). The text parameter contains the string in
// the input box.
func promptresponse(inp *prompt.Input, text string) {
	i3.Command(fmt.Sprintf("workspace %s", text))
}

// canceled is the callback that gets executed whenever the prompt is canceled.
// This can occur when the user presses escape (the "cancel" key).
func promptcanceled(inp *prompt.Input) {
	xevent.Quit(inp.X)
}

func addWorkspace() {
	// The input box uses the keybind module, so we must initialize it.
	keybind.Initialize(x)

	// Creating a new input prompt is as simple as supply an X connection,
	// a theme and a configuration. We use built in defaults here.
	inpPrompt := prompt.NewInput(x,
		prompt.DefaultInputTheme, prompt.DefaultInputConfig)

	// Show maps the input prompt window and sets the focus. It returns
	// immediately, and the main X event loop is started.
	// Also, we use the root window geometry to make sure the prompt is
	// centered in the middle of the screen. 'response' and 'canceled' are
	// callback functions.
	inpPrompt.Show(xwindow.RootGeometry(x),
		"New Workspace name: ", promptresponse, promptcanceled)

	xevent.Main(x)
}

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
	flag.Usage = usage
}

func main() {
	x, err = xgbutil.NewConn()
	if err != nil {
		log.Fatalln(err)
	}

	i3, err = i3ipc.GetIPCSocket()
	if err != nil {
		log.Fatalln(err)
	}

	flag.Parse()

	if flagAddworkspace {
		addWorkspace()
	}
}
