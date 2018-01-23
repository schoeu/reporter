package handdlers

import (
	"../autils"
	"log"
	"strings"
	"text/template"
)

func MarkdownMaker(path string) {
	if path == "" {
		log.Fatal("Template path is not exist.")
		return
	}

	t := template.New("Markdown")
	t, err := t.ParseFiles(path)
	autils.ErrHadle(err)
	// err = t.Execute(os.Stdout, data)
	// autils.ErrHadle(err)
}

func Array2Md(text []string) string {
	return strings.Join(text, "|")
}
