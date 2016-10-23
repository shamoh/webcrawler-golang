package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "webcrawler-golang",
	Short: "Find all libraries included into the found pages.",
	Long: `Application is able to invoke provided terms in Google Search
	and the parse all found pages then extract all scripts libraries included
	using <script> tag in page.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}