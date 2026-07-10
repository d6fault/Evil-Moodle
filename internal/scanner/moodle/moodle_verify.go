package scanner

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func VerifyMoodle(url string) {
	resp, err := HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	found := false
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "name":
					name = attr.Val
				case "content":
					content = attr.Val
				}
			}
			if name == "keywords" && strings.HasPrefix(content, "moodle,") {
				found = true
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if found {
		fmt.Printf("[+] %s is running Moodle\n", url)
	} else {
		fmt.Printf("[-] %s Seems to not be running Moodle\n", url)
		os.Exit(1)
	}
}
