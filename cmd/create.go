package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"hcgatewood/ghissue/lib"
)

var createHelpLong = `Create GitHub Issues from file.

Prints created Issue numbers to stdout.

The input file contains a repo target, followed by hyphen-separated issues.

The first line of each issue contains metadata, while all following lines
comprise the Issue body.

"""
repo_owner/repo_name
---
Title | Labels | Assignees
Body
---
Title | Labels | Assignees
Body
---
Title | Labels | Assignees
Body
"""
`

var createCmd = &cobra.Command{
	Use:                        "create issues.txt",
	Short:                      "Create Issues from file",
	Long:                       strings.ReplaceAll(createHelpLong, `"""`, "```"), // replace """ with ```
	Args:                       cobra.MinimumNArgs(1),
	Run:                        runCreate,
	SuggestionsMinimumDistance: 1,
}

var globalConfig lib.Config

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&globalConfig.DryRun, "dryrun", false, "don't actually create the Issues")
	createCmd.Flags().BoolVar(&globalConfig.Info, "info", false, "print more info about the Issues")
	createCmd.Flags().BoolVar(&globalConfig.Open, "open", false, "open browser to view new Issues")
	createCmd.Flags().BoolVar(&globalConfig.Byline, "byline", true, "append hcgatewood/ghissue byline to Issue body")
}

func runCreate(cmd *cobra.Command, args []string) {
	input, err := read(args[0])
	if err != nil {
		exit(cmd, err)
	}
	input = lib.TrimInput(input)

	cfg := &lib.Config{
		Token:  rootToken,
		DryRun: globalConfig.DryRun,
		Info:   globalConfig.Info,
		Open:   globalConfig.Open,
		Byline: globalConfig.Byline,
	}
	_, err = lib.Create(cfg, input)
	if err != nil {
		exit(cmd, err)
	}
}
