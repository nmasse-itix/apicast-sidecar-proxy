package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
    "strconv"
)

// The MyResponseWriter is a wrapper around the standard http.ResponseWriter
// We need it to retrieve the http status code of an http response
type MyResponseWriter struct {
	Underlying http.ResponseWriter
	Status     int
}

func (mrw *MyResponseWriter) Header() http.Header {
	return mrw.Underlying.Header()
}

func (mrw *MyResponseWriter) Write(b []byte) (int, error) {
	return mrw.Underlying.Write(b)
}

func (mrw *MyResponseWriter) WriteHeader(s int) {
	mrw.Status = s
	mrw.Underlying.WriteHeader(s)
}

func SetupReverseProxy(local_port int, target string) {
	// Parse the target URL and perform some sanity checks
	url, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	// Initialize a Reverse Proxy object with a custom director
	proxy := httputil.NewSingleHostReverseProxy(url)
	underlying_director := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Let the underlying director do the mandatory job
		underlying_director(req)

		// Custom Handling
		// ---------------
		//
		// Filter out the "Host" header sent by the client
		// otherwise the target server won't be able to find the
		// matching virtual host. The correct host header will be
		// added automatically by the net/http package.
		req.Host = ""
	}
    handler := http.NewServeMux()
	handler.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		// Log the incoming request (including headers)
		fmt.Printf("%v %v HTTP/1.1\n", req.Method, req.URL)
		req.Header.Write(os.Stdout)
		fmt.Println()

		// Wrap the standard response writer with our own
		// implementation because we need the status code of the
		// response and that field is not exported by default
		mrw := &MyResponseWriter{Underlying: rw}

		// Let the reverse proxy handle the request
		proxy.ServeHTTP(mrw, req)

		// Log the response
		fmt.Printf("%v %v\n", mrw.Status, http.StatusText(mrw.Status))
		mrw.Header().Write(os.Stdout)
		fmt.Println()
	})

	fmt.Printf("Listening on port %v for incoming requests...\n", local_port)
    err = http.ListenAndServe(fmt.Sprintf(":%v", local_port), handler)
    if (err != nil) {
      fmt.Println("ERROR: %s", err)
    }
}

func main() {
	portal_endpoint := os.Getenv("THREESCALE_PORTAL_ENDPOINT")
	backend_endpoint := os.Getenv("BACKEND_ENDPOINT_OVERRIDE")

    portal_port_opt := os.Getenv("PORTAL_LISTEN_PORT")
    if (portal_port_opt == "") {
        portal_port_opt = "9090"
        fmt.Println("WARNING: No PORTAL_LISTEN_PORT environment variable found, defaulting to '9090'...")
    }

    backend_port_opt := os.Getenv("BACKEND_LISTEN_PORT")
    if (backend_port_opt == "") {
        backend_port_opt = "9091"
        fmt.Println("WARNING: No BACKEND_LISTEN_PORT environment variable found, defaulting to '9091'...")
    }

	error := false
	if portal_endpoint == "" {
		fmt.Println("ERROR: No THREESCALE_PORTAL_ENDPOINT environment variable found !")
		error = true
	}
	if backend_endpoint == "" {
		fmt.Println("ERROR: No BACKEND_ENDPOINT_OVERRIDE environment variable found !")
		error = true
	}
    portal_port, err := strconv.Atoi(portal_port_opt)
    if (err != nil) {
        fmt.Printf("ERROR: Cannot parse the PORTAL_LISTEN_PORT environment variable (%s): %s\n", portal_port_opt, err)
		error = true
    }
    backend_port, err := strconv.Atoi(backend_port_opt)
    if (err != nil) {
        fmt.Printf("ERROR: Cannot parse the BACKEND_LISTEN_PORT environment variable (%s): %s\n", backend_port_opt, err)
		error = true
    }
	if error {
		os.Exit(1)
	}

	go SetupReverseProxy(backend_port, backend_endpoint)
	SetupReverseProxy(portal_port, portal_endpoint)
}
