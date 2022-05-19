package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/polyapi/pkg/proxy"
)

func main() {
	var maxIdle int
	var maxIdlePerHost int
	var namespace string
	flag.IntVar(&maxIdle, "maxIdle", 100, "max idle.")
	flag.IntVar(&maxIdlePerHost, "maxIdlePerHost", 100, "max idle per host.")
	flag.StringVar(&namespace, "namespace", "serving", "proxy namespace")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	transport := proxy.NewTransport(maxIdle, maxIdlePerHost)
	proxy, err := proxy.NewProxy(
		transport,
		namespace,
	)
	if err != nil {
		return
	}

	proxy.Run(ctx)
}
