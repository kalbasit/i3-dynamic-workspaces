package dmenu

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func Run(list []string) string {
	items := strings.Join(list, "\n")

	dmenu := exec.Command("rofi", "-dmenu")
	dmenu.Stdin = strings.NewReader(items)
	var out bytes.Buffer
	dmenu.Stdout = &out
	err := dmenu.Run()
	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}
