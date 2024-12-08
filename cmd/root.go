package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Web scraper",
	Short: "Find dead link",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		url := args[0]
		if len(url) == 0 {
			return
		}
		fmt.Printf("Scanning: %s\n", args[0])

		res, err := http.Get(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
