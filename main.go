package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/basgys/alterego/router"
)

var (
	rt                 *router.Router // Let's make it global
	requestLogging     bool
	redirectStatusCode int
)

func main() {
	// Basic configuration
	requestLogging = os.Getenv("REQUEST_LOGGING") == "true"
	rsc := os.Getenv("REDIRECT_STATUS_CODE")
	if rsc != "" {
		code, err := strconv.Atoi(rsc)
		if err != nil {
			fmt.Println("Cannot parse REDIRECT_STATUS_CODE", err)
			os.Exit(1)
		}
		redirectStatusCode = code
	} else {
		redirectStatusCode = http.StatusPermanentRedirect
	}

	// Build HTTP server addr
	ip := os.Getenv("IP")
	if ip == "" {
		ip = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5080"
	}
	addr := strings.Join([]string{ip, port}, ":")

	// Build router
	rules, err := buildRules()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	printRules(rules)

	rt = &router.Router{
		Rules: rules,
	}

	// Attach routes
	mux := http.NewServeMux()
	mux.HandleFunc("/__health__", health)
	mux.HandleFunc("/", redirect)

	// Boot
	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("\nRunning on %s...\n", addr)
	log.Fatal(s.ListenAndServe())
}

// health is called by the container health check
func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	uri, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req := url.URL{
		Scheme:   "http",
		Host:     r.Host,
		Path:     uri.Path,
		RawQuery: uri.RawQuery,
	}

	// Find matching route
	redir, ok := rt.Match(&req)
	if !ok {
		if requestLogging {
			fmt.Printf("[%s] No match %s -> %d\n",
				time.Now().UTC().Format(time.RFC3339Nano),
				req.String(),
				http.StatusNotFound,
			)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if requestLogging {
		fmt.Printf("[%s] Redirect %s -> %s\n",
			time.Now().UTC().Format(time.RFC3339Nano),
			req.String(),
			redir.String(),
		)
	}
	http.Redirect(w, r, redir.String(), http.StatusTemporaryRedirect)
}

func buildRules() ([]router.Rule, error) {
	var rules []router.Rule
	items := strings.Split(os.Getenv("REDIRECTS"), ",")
	for i, item := range items {
		if item == "" {
			return nil, fmt.Errorf("empty redirect #%d", i)
		}

		env := strings.TrimSpace(item)
		l := strings.Split(os.Getenv(env), ",")
		if len(l) != 2 {
			return nil, fmt.Errorf("malformed redirection %s from env var %s", l, env)
		}
		src, err := url.Parse(l[0])
		if err != nil {
			return nil, fmt.Errorf("bad source url %s", l[0])
		}
		dst, err := url.Parse(l[1])
		if err != nil {
			return nil, fmt.Errorf("bad destination url %s", l[1])
		}

		rules = append(rules, router.Rule{
			Src: *src,
			Dst: *dst,
		})
	}
	return rules, nil
}

func printRules(rules []router.Rule) {
	w := new(tabwriter.Writer)
	buf := bytes.NewBuffer([]byte{})
	w.Init(buf, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Source \t Destination")

	for _, r := range rules {
		fmt.Fprintf(w, "%s \t %s\n", r.Src.String(), r.Dst.String())
	}

	w.Flush()
	fmt.Println("Rules:")
	fmt.Println(buf.String())
}
