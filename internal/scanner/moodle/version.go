package scanner

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// looksLikeHTML reports whether a fetched body is an HTML document rather
// than the raw version artifact we asked for. A 404 or error page served
// with a 200 still contains version-like numbers, so we must reject it
// before parsing to avoid reporting a bogus version.
func looksLikeHTML(body string) bool {
	head := strings.ToLower(body)
	if i := strings.Index(head, "<"); i >= 0 {
		head = head[i:]
	}
	return strings.Contains(head, "<!doctype") ||
		strings.Contains(head, "<html") ||
		strings.Contains(head, "<head") ||
		strings.Contains(head, "<body")
}

// versionFromBody extracts a Moodle version string from the content of a
// version artifact. UPGRADING.md uses "## X.Y" branch headers; version.php
// carries a $release/$version assignment. We anchor on those markers rather
// than the first dotted number in the file, which otherwise picks up dates,
// tracker IDs and dependency versions.
func versionFromBody(content string) (string, bool) {
	// UPGRADING.md: the first "## X.Y" heading is the current branch.
	headingRe := regexp.MustCompile(`(?m)^##\s+(\d+\.\d+(?:\.\d+)?)\s*$`)
	if m := headingRe.FindStringSubmatch(content); m != nil {
		return m[1], true
	}

	// version.php: $release = '5.3.0'; (three-part) or $version = 2026100300
	releaseRe := regexp.MustCompile(`(?i)\$release\s*=\s*['"](\d+\.\d+(?:\.\d+)?)['"]`)
	if m := releaseRe.FindStringSubmatch(content); m != nil {
		return m[1], true
	}

	// version.php also exposes $branch = '503' when release is absent.
	branchRe := regexp.MustCompile(`(?i)\$branch\s*=\s*['"](\d{2,3})['"]`)
	if m := branchRe.FindStringSubmatch(content); m != nil {
		b := m[1]
		if len(b) >= 2 {
			return b[:len(b)-1] + "." + b[len(b)-1:], true
		}
	}

	return "", false
}

func CheckMoodleVersion(url string) {
	base := strings.TrimSuffix(url, "/")

	// /version.php is the canonical artifact but normally executes to an
	// empty body; /lib/UPGRADING.md is readable only on some installs.
	for _, path := range []string{"/lib/UPGRADING.md", "/version.php", "/lib/upgrade.txt"} {
		artifactURL := base + path
		resp, err := HTTPClient.Get(artifactURL)
		if err != nil {
			fmt.Printf("    Error: %v\n", err)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("    Error reading response: %v\n", err)
			continue
		}

		content := string(body)
		if looksLikeHTML(content) {
			continue
		}
		if version, ok := versionFromBody(content); ok {
			fmt.Printf("    [+] Version: %s (via %s)\n", version, path)
			PrintVersionStatus(version)
			return
		}
	}

	fmt.Printf("    [-] Could not find version (version artifacts are not web-readable)\n")
}

func GetLatestMoodleVersion() (string, error) {
	pageURL := "https://download.moodle.org/releases/latest/"

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", "moodleprobe")

	resp, err := HTTPClient.Do(req)
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

	// The releases page lists several branches; the newest stable is not
	// necessarily first, so take the maximum version rather than the first
	// tag that happens to lack a "+" suffix.
	latest := ""
	for _, m := range matches {
		version := normalizeVersion(m[1])
		if latest == "" || compareVersions(version, latest) > 0 {
			latest = version
		}
	}

	return latest, nil
}

func normalizeVersion(v string) string {
	parts := strings.Split(v, ".")
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	return strings.Join(parts, ".")
}

// numericPrefix returns the leading integer of a version segment. Segments
// such as "0rc1" or "1beta" are treated as their numeric part so that
// pre-release tags comparing equal to the release they precede sort no higher.
func numericPrefix(seg string) int {
	end := 0
	for end < len(seg) && seg[end] >= '0' && seg[end] <= '9' {
		end++
	}
	if end == 0 {
		return 0
	}
	n, _ := strconv.Atoi(seg[:end])
	return n
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
			an = numericPrefix(aParts[i])
		}
		if i < len(bParts) {
			bn = numericPrefix(bParts[i])
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
