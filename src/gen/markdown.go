package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/hcgatewood/ghissue/src/cmd"
)

//go:generate go run markdown.go
func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, ".")
	if err != nil {
		log.Fatal(err)
	}
}
