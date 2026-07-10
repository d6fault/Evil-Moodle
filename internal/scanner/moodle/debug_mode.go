package scanner

import (
	"fmt"
	"io"
	"strings"
)

func CheckDebugMode(url string) {
	fmt.Printf("[*] Checking debug mode\n")

	resp, err := HTTPClient.Get(url)
	if err != nil {
		fmt.Printf("    Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("    Error reading response: %v\n", err)
		return
	}
	content := string(body)

	switch {
	case strings.Contains(content, `"debug":false`):
		fmt.Printf("   [\x1b[31m-\x1b[0m] debug:false — debug mode is OFF\n")
	case strings.Contains(content, `"debug":true`):
		fmt.Printf("   [\x1b[32m+\x1b[0m] debug:true — debug mode is ON\n")
	default:
		fmt.Printf("    [?] No debug flag found — unknown\n")
	}
}
