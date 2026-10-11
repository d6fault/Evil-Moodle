package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	scanner "github.com/d6fault/evilmoodle/internal/scanner/moodle"
	"github.com/d6fault/evilmoodle/internal/scanner/ui"
	"github.com/spf13/cobra"
)

var url string
var wordlist string
var verbose bool
var rateLimit int
var delay int
var insecure bool

func scheme(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

var rootCmd = &cobra.Command{
	Use:   "evilmoodle",
	Short: "A Moodle reconnaissance and security scanning tool",
	Long: `Evil Moodle is a CLI tool for reconnaissance and security scanning of Moodle LMS instances.

It fingerprints the target, and enumerates installed
plugins via probing using a configurable wordlist.

Usage:
  evilmoodle --url <target> [flags]

Examples:
  evilmoodle --url https://moodle.example.com
  evilmoodle --url moodle.example.com --wordlist plugins.txt -v
  evilmoodle --url https://example.com --wordlist plugins.txt --rate 10 --delay 100`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()
		if url == "" {
			fmt.Println("Please Enter a url. With --url")
			os.Exit(1)
		}
		if insecure {
			scanner.SetInsecureTLS()
			fmt.Println("[!] TLS certificate verification disabled")
		}
		fmt.Println("Scanning", url)
		url = scheme(url)

		if !scanner.VerifyMoodle(url) {
			return
		}
		scanner.CheckMoodleVersion(url)
		scanner.CheckDebugMode(url)
		scanner.CheckSelfRegistration(url)
		scanner.CheckComposer(url)
		scanner.GuestLogin(url)

		if wordlist != "" {
			fmt.Println("\n[*] Starting plugin enumeration...")
			plugins, err := scanner.ProbePlugins(url, wordlist, verbose, rateLimit, time.Duration(delay)*time.Millisecond)
			if err != nil {
				fmt.Printf("[-] Failed to load wordlist: %v\n", err)
			} else {
				scanner.PrintFoundPlugins(plugins)
			}
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "Moodle URL")
	rootCmd.Flags().StringVar(&wordlist, "wordlist", "", "Path to plugins wordlist")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().IntVarP(&rateLimit, "rate", "r", 0, "Rate limit in requests per second (0 = unlimited)")
	rootCmd.Flags().IntVarP(&delay, "delay", "d", 0, "Delay between requests in milliseconds")
	rootCmd.Flags().BoolVar(&insecure, "insecure", false, "Skip TLS certificate verification")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
