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

	// Follow redirects for this client so we can inspect the page we
	// actually land on. Moodle redirects the login POST on both success
	// and failure, so the status code alone cannot distinguish them.
	client := &http.Client{
		Jar:     jar,
		Timeout: HTTPClient.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: HTTPClient.Transport,
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
	lowerBody := strings.ToLower(string(postBody))

	// A failed attempt re-renders the login form with an error and/or a
	// fresh logintoken. A successful attempt lands on a page carrying the
	// logged-in user menu instead.
	if strings.Contains(lowerBody, "invalidlogin") ||
		strings.Contains(lowerBody, "invalid login") ||
		strings.Contains(lowerBody, "loginerrormessage") ||
		strings.Contains(lowerBody, `name="logintoken"`) {
		fmt.Printf("    [-] Guest login failed\n")
		return
	}

	loggedIn := strings.Contains(lowerBody, "usermenu") ||
		strings.Contains(lowerBody, `href="/login/logout.php`) ||
		strings.Contains(lowerBody, `href="`+target+`"`) &&
			strings.Contains(lowerBody, "logout")

	if loggedIn {
		fmt.Printf("    [*] Guest login SUCCEEDED - guest:guest\n")
	} else {
		fmt.Printf("    [?] Guest login result inconclusive (status %d)\n", resp.StatusCode)
	}
}
