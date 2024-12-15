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

		var (
			visited     = make(map[string]struct{})
			results     = []Result{}
			fetchUrl    = make(chan string)
			originalUrl = args[0]
		)

		fetch := func() {
			for currentUrl := range fetchUrl {
				var err error
				if currentUrl != originalUrl {
					if currentUrl[0] != '/' {
						fmt.Printf("Out of scope: %s\n", currentUrl)
						continue
					}

					currentUrl, err = url.JoinPath(originalUrl, currentUrl)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}

				if _, ok := visited[currentUrl]; ok {
					fmt.Printf("Already visited: %s\n", currentUrl)
					continue
				}

				go getResult(currentUrl, &visited, &results, fetchUrl)
			}
		}
		go fetch()

		fetchUrl <- originalUrl

		wg.Wait()
		close(fetchUrl)

		fmt.Println("\nResults:")
		for _, result := range results {
			fmt.Printf("%d <- %s\n", result.status, result.url)
		}
	},
}

func getResult(currentUrl string, visited *map[string]struct{}, results *[]Result, fetchUrl chan<- string) {
	wg.Add(1)
	defer wg.Done()

	fmt.Printf("Scanning: %s\n", currentUrl)
	res, err := http.Get(currentUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	mut.Lock()
	(*visited)[currentUrl] = struct{}{}
	*results = append(*results, Result{url: currentUrl, status: res.StatusCode})
	mut.Unlock()

	if res.StatusCode != http.StatusOK {
		return
	}

	tokenizer := html.NewTokenizer(res.Body)
	newUrls := hrefs(tokenizer)
	for _, each := range newUrls {
		fetchUrl <- each
	}
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
