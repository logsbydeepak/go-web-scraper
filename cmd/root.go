package cmd

import (
	"fmt"
	"io"
	"net/url"
	"sync"

	"net/http"
	"os"

	"example.com/pkg/output"
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
			fetchUrl    = make(chan string)
			results     = []Result{}
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

				mut.Lock()
				if _, ok := visited[currentUrl]; ok {
					mut.Unlock()
					fmt.Printf("Already visited: %s\n", currentUrl)
					continue
				}
				visited[currentUrl] = struct{}{}
				mut.Unlock()

				go getResult(currentUrl, &results, fetchUrl)
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

func getResult(currentUrl string, results *[]Result, fetchUrl chan<- string) {
	wg.Add(1)
	defer wg.Done()

	fmt.Printf("Scanning: %s\n", currentUrl)
	res, err := http.Get(currentUrl)
	if err != nil {
		output.Error("Failed to get fetch")
		return
	}

	mut.Lock()
	*results = append(*results, Result{url: currentUrl, status: res.StatusCode})
	mut.Unlock()

	if res.StatusCode != http.StatusOK {
		return
	}

	tokenizer := html.NewTokenizer(res.Body)
	newUrls := getAllHrefUrl(tokenizer)
	for _, each := range newUrls {
		fetchUrl <- each
	}
}

func getAllHrefUrl(tokenizer *html.Tokenizer) (urls []string) {
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			} else {
				output.Error("Failed process token")
				break
			}
		}

		tag, hasAttr := tokenizer.TagName()
		if !(string(tag) == "a") && !hasAttr {
			continue
		}
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
