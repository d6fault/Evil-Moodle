package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

func CheckComposer(url string) {
	base := strings.TrimSuffix(url, "/")
	composerURL := base + "/composer.json"

	fmt.Printf("[*] Checking composer.json\n")

	resp, err := http.Get(composerURL)
	if err != nil {
		fmt.Printf("    Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("    [-] composer.json not found (status %d)\n", resp.StatusCode)
		return
	}

	fmt.Printf("    [+] composer.json is accessible\n")

	body, _ := io.ReadAll(resp.Body)

	// Parse the JSON — we only care about require-dev
	var composer struct {
		RequireDev map[string]string `json:"require-dev"`
	}

	if err := json.Unmarshal(body, &composer); err != nil {
		fmt.Printf("    Error parsing JSON: %v\n", err)
		return
	}

	if len(composer.RequireDev) == 0 {
		fmt.Printf("    [-] require-dev is empty or not present\n")
		return
	}

	fmt.Printf("    [+] require-dev packages found (%d):\n", len(composer.RequireDev))

	// Sort the keys for consistent output
	keys := make([]string, 0, len(composer.RequireDev))
	for k := range composer.RequireDev {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, pkg := range keys {
		version := composer.RequireDev[pkg]
		fmt.Printf("        %-40s %s\n", pkg, version)
	}
}
