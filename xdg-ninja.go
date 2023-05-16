package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

type Files struct {
	Path     string `json:"path"`
	Moveable bool   `json:"moveable"`
	Help     string `json:"help"`
}

type Programs struct {
	Name  string  `json:"name"`
	Files []Files `json:"files"`
}

//go:embed programs
var f embed.FS

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	entries, err := os.ReadDir(home)
	if err != nil {
		log.Fatal(err)
	}

	r, _ := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width (default is 80)
		glamour.WithWordWrap(80),
	)
	output := ""
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			input := fmt.Sprintf("programs/%s.json", e.Name()[1:])
			data, err := f.ReadFile(input)
			if err != nil {
				continue
			}

			var p Programs
			err = json.Unmarshal(data, &p)
			if err != nil {
				log.Fatalln(err)
			}

			name := fmt.Sprintf("# %s\n", p.Name)
			output += name

			for _, file := range p.Files {
				path := fmt.Sprintf("## %s\n", file.Path)
				output += path
				output += file.Help
			}
		}
	}

	o, err := r.Render(output)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(o)
}
