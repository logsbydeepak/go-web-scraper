package cmd

import (
	"fmt"
	"io"
	"net/url"

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
		if len(args[0]) == 0 {
			return
		}
		urls = append(urls, args[0])
		visited := make(map[string]struct{})

		for {
			var err error
			if len(urls) == 0 {
				break
			}

			currentUrl := urls[0]

			if currentUrl != args[0] {
				if currentUrl[0] != '/' {
					fmt.Printf("Out of scope: %s\n", currentUrl)
					urls = urls[1:]
					continue
				} else {
					currentUrl, err = url.JoinPath(args[0], currentUrl)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}

			_, ok := visited[currentUrl]
			if ok {
				fmt.Printf("Already visited: %s\n", currentUrl)
				urls = urls[1:]
				continue
			}

			fmt.Printf("Scanning: %s\n", currentUrl)
			res, err := getHtml(currentUrl)
			if err != nil {
				fmt.Println(err)
				return
			}

			tokenizer := html.NewTokenizer(res.Body)
			newUrls := hrefs(tokenizer)
			urls = append(urls, newUrls...)
			fmt.Println(urls)
			visited[currentUrl] = struct{}{}
			urls = urls[1:]
		}

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
