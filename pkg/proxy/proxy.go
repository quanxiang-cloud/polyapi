package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	port = 9999
)

// Proxy faas internal proxy
type Proxy interface {
	Run(ctx context.Context)
}

type proxy struct {
	server    *http.Server
	transport http.RoundTripper

	namespace string
}

// NewProxy return proxy
func NewProxy(transport http.RoundTripper, namespace string) (Proxy, error) {
	mux := http.NewServeMux()

	p := &proxy{
		namespace: namespace,
		transport: transport,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

	err := p.setHandle(mux)

	return p, err
}

func (p *proxy) setHandle(mux *http.ServeMux) error {
	mux.HandleFunc("/", p.do)
	return nil
}

func (p *proxy) do(w http.ResponseWriter, r *http.Request) {
	value := r.Header.Get("Refer")
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value = value[strings.LastIndex(value, "/")+1:]
	ok, err := regexp.MatchString(`^\w*-\w*-\w*\.r$`, value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value = value[:len(value)-2]
	values := strings.Split(value, "-")

	uri := fmt.Sprintf("http://%s-%s-%s-00001.%s.svc.cluster.local", values[2], values[1], values[0], p.namespace)

	url, err := url.ParseRequestURI(uri)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = p.transport
	r.Host = url.Host

	proxy.ServeHTTP(w, r)
}

func (p *proxy) Run(ctx context.Context) {
	doneCh := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				time.Second*5,
			)
			defer cancel()
			p.server.Shutdown(shutdownCtx) // nolint: errcheck
		case <-doneCh:
		}
	}()

	err := p.server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
	close(doneCh)
}
