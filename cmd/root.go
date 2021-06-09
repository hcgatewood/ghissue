package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flag vars
	rootTokenFilepath string
)

var rootCmd = &cobra.Command{
	Use:     "ghissue",
	Example: "ghissue --token ~/.creds/gh.token create issues.txt",
	Short:   "Bulk-upload GitHub Issues",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootTokenFilepath, "token", "", "Filepath of GitHub personal access token (required)")
	err := rootCmd.MarkPersistentFlagRequired("token")
	if err != nil {
		panic(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
