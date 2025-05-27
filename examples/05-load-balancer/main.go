package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"

	"github.com/borderzero/border0-go"
	"github.com/borderzero/border0-go/listen"
)

var servers = []string{
	"https://example.com",
	"https://example.org",
	"https://example.net",
}

var counter int32

func loadBalancerHandler(w http.ResponseWriter, r *http.Request) {
	// Get the next server URL.
	serverURL := getNextServer()

	// Parse the server URL.
	target, _ := url.Parse(serverURL)

	// Create a new reverse proxy.
	proxy := httputil.NewSingleHostReverseProxy(target)
	// Modify the request's host to the target's host.
	r.Host = target.Host

	// Remove the Accept-Encoding header to prevent gzip responses.
	r.Header.Del("Accept-Encoding")

	// Capture the headers from the incoming request.
	name := r.Header.Get("X-Auth-Name")
	email := r.Header.Get("X-Auth-email")

	// Remove the headers so they aren't sent upstream.
	r.Header.Del("X-Auth-Name")
	r.Header.Del("X-Auth-email")

	// Modify the response to insert text indicating the chosen backend server.

	proxy.ModifyResponse = func(resp *http.Response) error {

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// here we are modifying the content of the response
		// we'll let the user know what origin server was used example.com or example.org
		// Also insert the uers's name and email
		insertContent := fmt.Sprintf("<h1>Example Domain</h1><p style='color: red; font-weight: bold;'>Served by: %s<br>Welcome %s (%s)</p>", serverURL, name, email)
		body = []byte(strings.ReplaceAll(string(body), "<h1>Example Domain</h1>", insertContent))
		resp.Body = io.NopCloser(bytes.NewReader(body))
		return nil
	}

	// Serve the request using the reverse proxy.
	proxy.ServeHTTP(w, r)
}

func getNextServer() string {
	// Atomically increment the counter and get the next server.
	index := atomic.AddInt32(&counter, 1)
	return servers[index%int32(len(servers))]
}

func main() {
	listener, err := border0.Listen(
		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
	)
	if err != nil {
		log.Fatalln("failed to start listener:", err)
	}
	defer listener.Close()
	log.Println("Starting load balancer on Border0")
	// Use the border0 listener to serve the http handler (reverse proxy)
	log.Fatalln(http.Serve(listener, http.HandlerFunc(loadBalancerHandler)))

}
