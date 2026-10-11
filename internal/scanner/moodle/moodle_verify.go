package scanner

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// VerifyMoodle fetches the target and reports whether it looks like Moodle.
// It returns true when the target is Moodle and false otherwise; callers
// decide how to proceed so this can be used as a library function.
func VerifyMoodle(url string) bool {
	resp, err := HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return false
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

	// Fall back to body markers when the keywords meta tag is absent, which
	// is common on hardened or older installs.
	if !found {
		lower := strings.ToLower(string(body))
		if strings.Contains(lower, "moodle") &&
			(strings.Contains(lower, "/theme/") || strings.Contains(lower, "moodle.org") ||
				strings.Contains(lower, "mform") || strings.Contains(lower, "moodlelogo")) {
			found = true
		}
	}

	if found {
		fmt.Printf("[+] %s is running Moodle\n", url)
	} else {
		fmt.Printf("[-] %s Seems to not be running Moodle\n", url)
	}
	return found
}
