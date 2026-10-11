package scanner

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var insecureTLS bool

var HTTPClient = newClient()

func newClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Jar:     mustJar(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureTLS,
			},
		},
	}
}

// SetInsecureTLS rebuilds the shared client so that TLS certificate
// verification is skipped. It must be called before any scan starts.
func SetInsecureTLS() {
	insecureTLS = true
	HTTPClient = newClient()
}

func mustJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	return jar
}
