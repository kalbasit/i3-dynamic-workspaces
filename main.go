package main

import (
	"bytes"

	"flag"

	"image/color"

	"fmt"

	"io/ioutil"

	"github.com/proxypoke/i3ipc"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"

	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/wingo/render"

	"log"

	"os"

	"path"

	"strings"
)

var font = xgraphics.MustFont(xgraphics.ParseFont(
	bytes.NewBuffer(readFile("/usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf"))))

var defaultInputTheme = &prompt.InputTheme{
	BorderSize:  5,
	BgColor:     render.NewImageColor(color.RGBA{0xff, 0xff, 0xff, 0xff}),
	BorderColor: render.NewImageColor(color.RGBA{0x0, 0x0, 0x0, 0xff}),
	Padding:     10,

	Font: font,

	FontSize:   20.0,
	FontColor:  render.NewImageColor(color.RGBA{0x0, 0x0, 0x0, 0xff}),
	InputWidth: 500,
}

var defaultSelectTheme = &prompt.SelectTheme{
	BorderSize:  10,
	BgColor:     render.NewImageColor(color.RGBA{0xff, 0xff, 0xff, 0xff}),
	BorderColor: render.NewImageColor(color.RGBA{0x0, 0x0, 0x0, 0xff}),
	Padding:     20,

	Font: font,

	FontSize:  20.0,
	FontColor: render.NewImageColor(color.RGBA{0x0, 0x0, 0x0, 0xff}),

	ActiveBgColor:   render.NewImageColor(color.RGBA{0x0, 0x0, 0x0, 0xff}),
	ActiveFontColor: render.NewImageColor(color.RGBA{0xff, 0xff, 0xff, 0xff}),

	GroupBgColor: render.NewImageColor(color.RGBA{0xff, 0xff, 0xff, 0xff}),
	GroupFont:    font,

	GroupFontSize:  25.0,
	GroupFontColor: render.NewImageColor(color.RGBA{0x33, 0x66, 0xff, 0xff}),
	GroupSpacing:   15,
}

type item struct {
	text       string
	workspace  *i3ipc.Workspace
	promptItem *prompt.SelectItem
}

func newItem(workspace *i3ipc.Workspace) *item {
	return &item{
		text:       workspace.Name,
		workspace:  workspace,
		promptItem: nil,
	}
}

func readFile(fpath string) []byte {
	content, err := ioutil.ReadFile(fpath)

	if err != nil {
		log.Fatalln(err)
	}
	return content
}

func (item *item) SelectText() string {
	return item.text
}

func (item *item) SelectHighlighted(data interface{}) {}

func (item *item) SelectSelected(data interface{}) {
	i3.Command(fmt.Sprintf("workspace %s", item.text))
	os.Exit(0)
}

var (
	x   *xgbutil.XUtil
	i3  *i3ipc.IPCSocket
	err error

	flagAddworkspace    bool
	flagSwitchworkspace bool
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
		defaultInputTheme, prompt.DefaultInputConfig)

	// Show maps the input prompt window and sets the focus. It returns
	// immediately, and the main X event loop is started.
	// Also, we use the root window geometry to make sure the prompt is
	// centered in the middle of the screen. 'response' and 'canceled' are
	// callback functions.
	inpPrompt.Show(xwindow.RootGeometry(x),
		"New Workspace name: ", promptresponse, promptcanceled)

	xevent.Main(x)
}

func switchWorkspace() {
	// The input box uses the keybind module, so we must initialize it.
	keybind.Initialize(x)

	slct := prompt.NewSelect(x,
		defaultSelectTheme, prompt.DefaultSelectConfig)

	workspaces, err := i3.GetWorkspaces()

	if err != nil {
		log.Fatalln(err)
	}

	items := make([]*prompt.SelectItem, 0, len(workspaces))

	for _, workspace := range workspaces {
		item := newItem(&workspace)
		item.promptItem = slct.AddChoice(item)

		items = append(items, item.promptItem)
	}

	group := slct.AddGroup(slct.NewStaticGroup("Workspaces"))

	groups := make([]*prompt.SelectShowGroup, 0, 1)
	groups = append(groups, group.ShowGroup(items))

	slct.Show(xwindow.RootGeometry(x), prompt.TabCompletePrefix, groups, slct)

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
	flag.BoolVar(&flagSwitchworkspace, "switch-workspace", false, "Switch to an existing workspace")

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
	} else if flagSwitchworkspace {
		switchWorkspace()
	}
}
