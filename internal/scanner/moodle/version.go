package scanner

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func CheckMoodleVersion(url string) {
	base := strings.TrimSuffix(url, "/")

	upgradeTxtURL := base + "/lib/upgrade.txt"
	resp, err := http.Get(upgradeTxtURL)
	if err != nil {
		fmt.Printf("    Error: %v\n", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("    Status: %d — file not found\n", resp.StatusCode)
		resp.Body.Close()
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if !strings.Contains(string(body), "4.5 Onwards") {
		lines := strings.Split(string(body), "\n")
		versionRegex := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)
		for _, line := range lines {
			if match := versionRegex.FindStringSubmatch(line); match != nil {
				fmt.Printf("    [+] Version: %s\n", match[1])
				PrintVersionStatus(match[1])
				return
			}
		}
		fmt.Printf("    [-] Could not find version\n")
		return
	}

	upgradingMdURL := base + "/lib/UPGRADING.md"
	fmt.Printf("[*] Checking %s\n", upgradingMdURL)

	resp2, err := http.Get(upgradingMdURL)
	if err != nil {
		fmt.Printf("    Error: %v\n", err)
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		fmt.Printf("    Status: %d — UPGRADING.md not found\n", resp2.StatusCode)
		return
	}

	body2, _ := io.ReadAll(resp2.Body)
	content := string(body2)

	lines := strings.Split(content, "\n")
	versionRegex := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)

	for i, line := range lines {
		if i >= 5 {
			break
		}
		if match := versionRegex.FindStringSubmatch(line); match != nil {
			fmt.Printf("    [+] Version: %s\n", match[1])
			PrintVersionStatus(match[1])
			return
		}
	}

	fmt.Printf("    [-] Could not find version\n")
}

func GetLatestMoodleVersion() (string, error) {
	pageURL := "https://download.moodle.org/releases/latest/"

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", "moodleprobe")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download.moodle.org returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	html := string(body)

	versionTagRegex := regexp.MustCompile(`<strong>\s*Moodle\s+(\d+\.\d+(?:\.\d+)?)(\+)?\s*</strong>`)

	matches := versionTagRegex.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("no version tags found on page")
	}

	for _, m := range matches {
		version := m[1]    // e.g. "5.2.1"
		plusSuffix := m[2] // e.g. "+" or ""

		if plusSuffix != "+" {
			return version, nil
		}
	}

	if len(matches) > 0 {
		return strings.TrimSuffix(matches[0][1], "+"), nil
	}

	return "", fmt.Errorf("could not parse version from page")
}

func normalizeVersion(v string) string {
	parts := strings.Split(v, ".")
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	return strings.Join(parts, ".")
}

func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var an, bn int
		if i < len(aParts) {
			fmt.Sscanf(aParts[i], "%d", &an)
		}
		if i < len(bParts) {
			fmt.Sscanf(bParts[i], "%d", &bn)
		}
		if an != bn {
			return an - bn
		}
	}
	return 0
}

func PrintVersionStatus(targetVersion string) {
	latest, err := GetLatestMoodleVersion()
	if err != nil {
		fmt.Printf("    [!] Could not fetch latest version: %v\n", err)
		return
	}

	targetNorm := normalizeVersion(targetVersion)
	latestNorm := normalizeVersion(latest)
	cmp := compareVersions(targetNorm, latestNorm)

	switch {
	case cmp < 0:
		fmt.Printf("    [*] OUTDATED — target: %s, latest: %s\n", targetNorm, latestNorm)
	case cmp == 0:
		fmt.Printf("    [*] UP TO DATE — %s\n", targetNorm)
	default:
		fmt.Printf("    [?] Target (%s) is newer than latest known (%s) — check version format\n",
			targetNorm, latestNorm)
	}
}
