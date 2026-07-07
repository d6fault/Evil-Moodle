package scanner

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CheckDebugMode(url string) {
	fmt.Printf("[*] Checking debug mode\n")

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("    Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	content := string(body)

	switch {
	case strings.Contains(content, `"debug":false`):
		fmt.Printf("    [-] debug:false — debug mode is OFF\n")
	case strings.Contains(content, `"debug":true`):
		fmt.Printf("    [+] debug:true — debug mode is ON\n")
	default:
		fmt.Printf("    [?] No debug flag found — unknown\n")
	}
}
