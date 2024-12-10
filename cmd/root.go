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

		var urls []string

		url := args[0]
		if len(url) == 0 {
			return
		}
		urls = append(urls, args[0])

		fmt.Printf("Scanning: %s\n", args[0])

		for _, url := range urls {
			res, err := getHtml(url)
			if err != nil {
				fmt.Println(err)
				return
			}

			tokenizer := html.NewTokenizer(res.Body)
			newUrls := hrefs(tokenizer)
			urls = append(urls, newUrls...)
		}
		fmt.Println(urls)
	},
}

func getHtml(url string) (*http.Response, error) {
	return http.Get(url)
}

func hrefs(tokenizer *html.Tokenizer) (urls []string) {
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			} else {
				fmt.Println(tokenizer.Err())
				break
			}
		}

		tag, hasAttr := tokenizer.TagName()
		if string(tag) == "a" && hasAttr {
			for {
				attrKey, attrValue, moreAttr := tokenizer.TagAttr()
				if string(attrKey) == "href" {
					urls = append(urls, string(attrValue))
				}
				if !moreAttr {
					break
				}
			}
		}
	}

	return urls
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
