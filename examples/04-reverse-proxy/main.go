package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
)

// This is the code we inject.
// Just for demonstration purposes, we're replacing the title of the page with a blinking text.
const appendScript = `
<style>
  .blink {
    animation: blinker 1s linear infinite;
  }
  @keyframes blinker {
    50% {
      opacity: 0;
    }
  }
</style>
<script>
	const replaceTitleInterval = setInterval(() => {
		const title = document.querySelector('.module--header h2 span');
		if (title) {
			title.innerHTML = '<span class="blink">Welcome to Border0</span><br><marquee bgcolor="Green" direction="left" >Your own version of the BBC</marquee>'
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

		resp.Body = io.NopCloser(bytes.NewReader(body))
		return nil
	}

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// border0 will set this header along with a few other identity related headers
			// now we know the users name, we can use it to personalize the page
			name := r.Header.Get("X-Auth-Name")
			email := r.Header.Get("X-Auth-email")
			log.Println("serving request from", name, email)
			// serve the reverse proxy with the correct host header
			r.Host = "www.bbc.com"
			proxy.ServeHTTP(w, r)
		},
	)

	// Use the border0 listener to serve the http handler (reverse proxy)
	log.Fatalln(http.Serve(listener, handler))
}
