package scanner

import "testing"

func TestVersionFromBody(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
		ok      bool
	}{
		{
			name:    "upgrading md heading",
			content: "# Moodle Upgrade notes\n\nIntro text with 2025.01 dates.\n\n## 5.3\n\n### core\n",
			want:    "5.3",
			ok:      true,
		},
		{
			name:    "version php release",
			content: "<?php\n$release = '4.5.1';\n$version = 2025010100;\n",
			want:    "4.5.1",
			ok:      true,
		},
		{
			name:    "version php branch",
			content: "<?php\n$branch = '503';\n$maturity = MATURITY_STABLE;\n",
			want:    "50.3",
			ok:      true,
		},
		{
			name:    "no version markers",
			content: "some text with a date 2025-01-01 and number 12345\n",
			want:    "",
			ok:      false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := versionFromBody(tc.content)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if got != tc.want {
				t.Fatalf("version = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLooksLikeHTML(t *testing.T) {
	if !looksLikeHTML("<html><head><title>404</title></head></html>") {
		t.Fatal("expected HTML to be detected")
	}
	if looksLikeHTML("<?php\n$release = '4.5.1';") {
		t.Fatal("PHP source misdetected as HTML")
	}
}

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"4.5.1", "5.3.0", -1},
		{"5.3.0", "5.3", 0},
		{"5.10", "5.9", 1},
		{"5.3.0", "5.3.0", 0},
	}
	for _, tc := range cases {
		got := compareVersions(tc.a, tc.b)
		switch {
		case got < 0 && tc.want >= 0, got > 0 && tc.want <= 0, got == 0 && tc.want != 0:
			t.Fatalf("compareVersions(%q,%q) = %d, want sign %d", tc.a, tc.b, got, tc.want)
		}
	}
}
