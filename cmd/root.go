package cmd

import (
	"fmt"
	"io"
	"net/url"
	"sync"

	"net/http"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

type Result struct {
	url    string
	status int
}

var wg sync.WaitGroup
var mut sync.Mutex

var rootCmd = &cobra.Command{
	Use:   "Web scraper",
	Short: "Find dead link",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		if len(args[0]) == 0 {
			return
		}

		visited := make(map[string]struct{})
		var results []Result

		fetchUrl := make(chan string)

		go func() {
			for myurl := range fetchUrl {
				go func() {
					wg.Add(1)
					defer wg.Done()
					var err error
					currentUrl := myurl

					if currentUrl != args[0] {
						if currentUrl[0] != '/' {
							fmt.Printf("Out of scope: %s\n", currentUrl)
							return
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
						return
					}

					fmt.Printf("Scanning: %s\n", currentUrl)
					res, err := getHtml(currentUrl)
					if err != nil {
						fmt.Println(err)
						return
					}
					visited[currentUrl] = struct{}{}

					fmt.Println(res.Status)
					results = append(results, Result{url: currentUrl, status: res.StatusCode})
					if res.StatusCode == http.StatusOK {
						tokenizer := html.NewTokenizer(res.Body)
						newUrls := hrefs(tokenizer)
						for _, each := range newUrls {
							fetchUrl <- each
						}

					}
				}()
			}
		}()
		fetchUrl <- args[0]

		wg.Wait()
		close(fetchUrl)

		fmt.Println("\nResults:")
		for _, result := range results {
			fmt.Printf("%d <- %s\n", result.status, result.url)
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
