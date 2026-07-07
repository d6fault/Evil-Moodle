package cmd

import (
	"fmt"
	"os"
	"strings"

	scanner "github.com/d6fault/moodleprobe/internal/scanner/moodle"
	"github.com/d6fault/moodleprobe/internal/scanner/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "moodleprobe",
}

var url string

func Scheme(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

var ScanCMD = &cobra.Command{
	Use:   "scan",
	Short: "The Moodle App you would like to scan",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Print_Banner()
		if url == "" {
			fmt.Println("Please Enter a url. With --url")
			os.Exit(1)
		}
		fmt.Println("Scanning", url)
		url = Scheme(url)
		scanner.Moodle_test(url)
		scanner.CheckMoodleVersion(url)
		scanner.CheckDebugMode(url)
		scanner.CheckSelfRegistration(url)
		scanner.CheckComposer(url)
		scanner.GuestLogin(url)
	},
}

func init() {
	rootCmd.AddCommand()
	ScanCMD.Flags().StringVar(&url, "url", "", "Moodle URL")
	rootCmd.AddCommand(ScanCMD)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
