package cmd

import (
	"github.com/spf13/cobra"
	"github.com/petrbouda/webcrawler-golang/crawler"
)

var (
	crawler = gcrawler.NewCrawler()
	debug *bool
)

var storeCmd = &cobra.Command{
	Use:   "search",
	Short: "Search a term using Google and find all <script> tags urls in the found results.",
	Example: "search something something-else",
	Run: func(cmd *cobra.Command, terms []string) {
		if len(terms) > 1 {
			println("Search command requires one argument (a term to find using a Google Search)")
			return
		}

		crawler.Crawl(terms[0])
	},
}

func init() {
	debug = storeCmd.Flags().BoolP("debug", "d", false, "Enable verbose level in HTTP Client.")

	RootCmd.AddCommand(storeCmd)
}