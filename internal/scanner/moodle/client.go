package scanner

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var HTTPClient = &http.Client{
	Timeout: 15 * time.Second,
	Jar:     mustJar(),
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func mustJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	return jar
}
