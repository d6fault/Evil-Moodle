package scanner

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type FoundPlugin struct {
	Path       string
	StatusCode int
}

func ProbePlugins(baseURL string, wordlistPath string, verbose bool, rateLimit int, delay time.Duration) ([]FoundPlugin, error) {
	plugins, err := loadWordlist(wordlistPath)
	if err != nil {
		return nil, err
	}

	baseURL = strings.TrimRight(baseURL, "/")

	totalPlugins := len(plugins)
	fmt.Printf("[*] Loaded %d plugins from wordlist\n", totalPlugins)
	fmt.Printf("[*] Target: %s\n", baseURL)
	fmt.Printf("[*] Mode: sequential\n")
	if rateLimit > 0 {
		fmt.Printf("[*] Rate limit: %d req/s\n", rateLimit)
	}
	if delay > 0 {
		fmt.Printf("[*] Delay: %s between requests\n", delay)
	}
	fmt.Println(strings.Repeat("-", 60))

	client := *HTTPClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	var found []FoundPlugin

	var checked int64
	var notFound int64
	var errors int64
	var redirects int64

	var limiter <-chan time.Time
	if rateLimit > 0 {
		interval := time.Second / time.Duration(rateLimit)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		limiter = ticker.C
	}

	startTime := time.Now()

	for i, path := range plugins {
		if limiter != nil {
			<-limiter
		}

		url := fmt.Sprintf("%s/%s/version.php", baseURL, path)
		resp, err := client.Get(url)
		atomic.AddInt64(&checked, 1)

		current := atomic.LoadInt64(&checked)

		if err != nil {
			atomic.AddInt64(&errors, 1)
			if verbose {
				fmt.Printf("\033[33m[ERR]\033[0m [%d/%d] %s - %v\n", i+1, totalPlugins, path, err)
			}
			continue
		}

		switch resp.StatusCode {
		case 200:
			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body := string(buf[:n])

			if !strings.Contains(body, "<html") && !strings.Contains(body, "<!DOCTYPE") {
				found = append(found, FoundPlugin{
					Path:       path,
					StatusCode: resp.StatusCode,
				})
				fmt.Printf("\033[32m[+]\033[0m [%d/%d] FOUND: %s (200)\n", i+1, totalPlugins, path)
			} else {
				atomic.AddInt64(&redirects, 1)
				if verbose {
					fmt.Printf("\033[34m[302]\033[0m [%d/%d] %s - login redirect\n", i+1, totalPlugins, path)
				}
			}

		case 403:
			atomic.AddInt64(&redirects, 1)
			if verbose {
				fmt.Printf("\033[33m[403]\033[0m [%d/%d] %s - forbidden (may exist)\n", i+1, totalPlugins, path)
			}

		case 404:
			atomic.AddInt64(&notFound, 1)
			if verbose {
				fmt.Printf("\033[31m[404]\033[0m [%d/%d] %s\n", i+1, totalPlugins, path)
			}

		default:
			if verbose {
				fmt.Printf("\033[33m[%d]\033[0m [%d/%d] %s\n", resp.StatusCode, i+1, totalPlugins, path)
			}
		}

		resp.Body.Close()

		if delay > 0 {
			time.Sleep(delay)
		}

		if current%100 == 0 && current > 0 {
			elapsed := time.Since(startTime).Seconds()
			rate := float64(current) / elapsed
			remaining := float64(totalPlugins-int(current)) / rate
			fmt.Printf("\n[*] Progress: %d/%d (%.1f%%) | Rate: %.1f/s | ETA: %.0fs | Found: %d\n\n",
				current, totalPlugins, float64(current)/float64(totalPlugins)*100,
				rate, remaining, len(found))
		}
	}

	elapsed := time.Since(startTime)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("[*] Scan complete in %s\n", elapsed.Round(time.Second))
	fmt.Printf("[*] Checked: %d | Found: %d | 404: %d | Redirects: %d | Errors: %d\n",
		checked, len(found), notFound, redirects, errors)

	return found, nil
}

func loadWordlist(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var plugins []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			plugins = append(plugins, line)
		}
	}

	return plugins, sc.Err()
}

func PrintFoundPlugins(plugins []FoundPlugin) {
	if len(plugins) == 0 {
		fmt.Println("[-] No third-party plugins found via probing")
		return
	}

	fmt.Printf("\n[+] Found %d installed plugins:\n", len(plugins))
	for _, p := range plugins {
		fmt.Printf("  %s\n", p.Path)
	}
}
