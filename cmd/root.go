package cmd

import (
	"fmt"
	"io"
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

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return
		}
		bodyString := string(body)
		fmt.Println("Response Body:", bodyString)
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
