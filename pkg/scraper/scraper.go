package scraper

import (
	"io"

	"example.com/pkg/output"
	"golang.org/x/net/html"
)

func GetAllHrefUrl(tokenizer *html.Tokenizer) (urls []string) {
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
