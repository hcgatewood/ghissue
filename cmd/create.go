package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"hcgatewood/ghissue/lib"
)

var createHelpLong = `Create GitHub Issues from file.

Prints created Issue numbers to stdout.

File format:

repo_owner/repo_name
---
Issue title 0 | labelX,labelY | assigneeA,assigneeB
Body can cross
Multiple lines
And the triple-hyphen is the divider
---
Issue title 1 | labelZ | assigneeC
Body can also be empty
As well as labels and assignees
But title is required
---
Smallest possible issue (just the title)
`

var createCmd = &cobra.Command{
	Use:                        "create issues.txt",
	Short:                      "Create Issues from file",
	Long:                       createHelpLong,
	Args:                       cobra.MinimumNArgs(1),
	Run:                        runCreate,
	SuggestionsMinimumDistance: 1,
}

var globalConfig lib.Config

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&globalConfig.DryRun, "dryrun", false, "Don't actually create the Issues")
	createCmd.Flags().BoolVar(&globalConfig.Info, "info", false, "Print more info about the Issues")
	createCmd.Flags().BoolVar(&globalConfig.Open, "open", false, "Open browser to view new Issues")
}

// hcg "/Users/hcgatewood/Desktop/tmp/hcgatewood23.token"

func runCreate(cmd *cobra.Command, args []string) {
	cfg := &lib.Config{
		Token:  read(rootTokenFilepath),
		DryRun: globalConfig.DryRun,
		Info:   globalConfig.Info,
		Open:   globalConfig.Open,
	}
	_, err := lib.Create(cfg, args[0])
	if err != nil {
		_ = cmd.Usage()
		fmt.Println("")
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func read(filepath string) string {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("could not read from %s: %+v", filepath, err)
	}
	return lib.TrimInput(string(bytes))
}
