package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
)

const appendScript = `
<script>
	const replaceTitleInterval = setInterval(() => {
		const title = document.querySelector('.module--header h2 span');
		if (title) {
			title.innerHTML = 'Welcome to Border0'
			clearInterval(replaceTitleInterval);
		}
	}, 100)
</script>
`

func main() {
	listener, err := border0.Listen(
		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}

	// create a http single host reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "www.bbc.com",
	})
	proxy.ModifyResponse = func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// replace all occurrences of "BBC" with "Border0"
		body = []byte(strings.ReplaceAll(string(body), "BBC", "Border0") + appendScript)
		resp.Body = ioutil.NopCloser(bytes.NewReader(body))
		return nil
	}

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// serve the reverse proxy with the correct host header
			r.Host = "www.bbc.com"
			proxy.ServeHTTP(w, r)
		},
	)

	log.Fatalln(http.Serve(listener, handler))
}
