package scanner

import (
	"fmt"
	"net/http"
	"strings"
)

func CheckSelfRegistration(url string) {
	base := strings.TrimSuffix(url, "/")
	signupURL := base + "/login/signup.php"

	fmt.Printf("[*] Checking self-registration\n")

	resp, err := http.Get(signupURL)
	if err != nil {
		fmt.Printf("    Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Printf("    [-] Self-registration is CLOSED\n")
	} else if resp.StatusCode == 200 {
		fmt.Printf("    [+] Self-registration is OPEN\n")
	} else {
		fmt.Printf("    [?] Status %d — unknown\n", resp.StatusCode)
	}
}
