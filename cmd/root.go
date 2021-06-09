package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// Global flag vars
	rootTokenFilepath string
	rootToken         string
)

var (
	tokenEnvName = "GITHUB_TOKEN"
)

// RootCmd is the root command for the CLI.
var RootCmd = &cobra.Command{
	Use:              "ghissue",
	Example:          "GITHUB_TOKEN='...' ghissue create ./issues.txt",
	Short:            "Bulk-upload GitHub Issues",
	PersistentPreRun: globalPre,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&rootToken, "token", "", "GitHub personal access token")
}

// Execute the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Unexpected error:", err)
		os.Exit(1)
	}
}

func globalPre(cmd *cobra.Command, args []string) {
	if rootToken == "" {
		envToken := os.Getenv(tokenEnvName)
		if envToken == "" {
			exit(cmd, errors.New("Must provide a GitHub personal access token. Try setting GITHUB_TOKEN environment variable."))
		}
		rootToken = envToken
	}
	rootToken = strings.TrimSpace(rootToken)
}

func exit(cmd *cobra.Command, err error) {
	_ = cmd.Usage()
	fmt.Println("")
	fmt.Println("Error:", err)
	os.Exit(1)
}

func read(filepath string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", errors.Errorf("could not read from %s: %+v", filepath, err)
	}
	return string(bytes), nil
}
