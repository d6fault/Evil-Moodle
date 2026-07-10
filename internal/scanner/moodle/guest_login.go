package scanner

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

func GuestLogin(target string) {
	target = strings.TrimRight(target, "/")
	loginURL := target + "/login/index.php"

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Printf("[-] Guest login: failed to create cookie jar: %v\n", err)
		return
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: HTTPClient.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	fmt.Printf("[*] Checking Guest login\n")
	resp, err := client.Get(loginURL)
	if err != nil {
		fmt.Printf("[-] Guest login: request failed: %v\n", err)
		return
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("[-] Guest login: failed to read response: %v\n", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("[-] Guest login: login page returned status %d\n", resp.StatusCode)
		return
	}

	tokenRe := regexp.MustCompile(`<input[^>]*name="logintoken"[^>]*value="([^"]+)"`)
	match := tokenRe.FindStringSubmatch(string(body))
	token := ""
	if len(match) >= 2 {
		token = match[1]
	}

	form := url.Values{}
	form.Set("username", "guest")
	form.Set("password", "guest")
	if token != "" {
		form.Set("logintoken", token)
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Printf("[-] Guest login: failed to build request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "moodleprobe")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("[-] Guest login: POST failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	postBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[-] Guest login: failed to read response: %v\n", err)
		return
	}
	bodyStr := string(postBody)
	success := false

	if resp.StatusCode == 303 || resp.StatusCode == 302 || resp.StatusCode == 301 {
		success = true
	}

	if !success {
		lowerBody := strings.ToLower(bodyStr)
		if strings.Contains(lowerBody, "invalidlogin") || strings.Contains(lowerBody, "invalid login") {
			success = false
		}
	}

	if success {
		fmt.Printf("    [*] Guest login SUCCEEDED - guest:guest\n")
	} else {
		fmt.Printf("    [-] Guest login failed\n")
	}
}
