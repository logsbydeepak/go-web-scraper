package cmd

import (
	"fmt"
	"io"
	// "io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
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

		t := html.NewTokenizer(res.Body)

		for {
			tt := t.Next()
			if tt == html.ErrorToken {
				if t.Err() == io.EOF {
					break
				} else {
					fmt.Println(t.Err())
					break
				}
			}

			tag, hasAttr := t.TagName()
			if string(tag) == "a" && hasAttr {
				for {
					attrKey, attrValue, moreAttr := t.TagAttr()
					if string(attrKey) == "href" {
						fmt.Println(string(attrKey), string(attrValue))
					}
					if !moreAttr {
						break
					}
				}
			}

		}
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
