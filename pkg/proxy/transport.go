package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	sleep = 30 * time.Millisecond
)

// NewTransport return http round tripper
func NewTransport(maxIdle, maxIdlePerHost int) http.RoundTripper {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		bo := wait.Backoff{
			Duration: 50 * time.Millisecond,
			Factor:   1.4,
			Jitter:   0.1,
			Steps:    15,
		}

		dialer := &net.Dialer{
			Timeout:   bo.Duration, // Initial duration.
			KeepAlive: 5 * time.Second,
			DualStack: true,
		}

		start := time.Now()
		for {
			c, err := dialer.DialContext(ctx, network, address)
			if err != nil {
				var errNet net.Error
				if errors.As(err, &errNet) && errNet.Timeout() {
					if bo.Steps < 1 {
						break
					}
					dialer.Timeout = bo.Step()
					time.Sleep(wait.Jitter(sleep, 1.0)) // Sleep with jitter.
					continue
				}
				return nil, err
			}
			return c, nil
		}
		elapsed := time.Since(start)
		return nil, fmt.Errorf("timed out dialing after %.2fs", elapsed.Seconds())
	}

	transport.DisableKeepAlives = true
	transport.MaxIdleConns = maxIdle
	transport.MaxIdleConnsPerHost = maxIdlePerHost
	transport.ForceAttemptHTTP2 = false
	transport.DisableCompression = false

	return transport
}
