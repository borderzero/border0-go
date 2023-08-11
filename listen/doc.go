// Package listen can create a Border0 listener with configurable options. The border0 listener
// is a net.Listener that can be used to accept incoming connections. When the listener is passed to
// http.Serve, the server will accept HTTP requests sent by Border0 and forward them to an HTTP handler.
// The handler's response will be sent back to Border0. If no options are provided, some default values
// will be used.
//
// Example:
//
//	listener := listen.New(
//		listen.WithSocketName("sdk-socket-http"),              // http socket name the listener will be bound to, socket will be created if not exists
//		listen.WithAuthToken(os.Getenv("BORDER0_AUTH_TOKEN")), // optional, if not provided, Border0 SDK will use BORDER0_AUTH_TOKEN env var
//	)
//	err := listener.Start()
//	if err != nil {
//		// handle error
//	}
//	defer listener.Close()
//
//	// create a handler
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// do something
//	})
//
//	// start a server using the listener
//	http.Serve(listener, handler)
//
// See [Option] for more configurable options.
package listen
