package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"hcgatewood/ghissue/cmd"
)

//go:generate go run markdown.go
func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, ".")
	if err != nil {
		log.Fatal(err)
	}
}
