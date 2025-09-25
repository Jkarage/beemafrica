package beemafrica

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"
)

const version = "v1"

var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

type Logger func(context.Context, string, ...any)

type Client struct {
	log       Logger
	apiKey    string
	apiSecret string
	http      *http.Client
}

func generateBasicHeader(apiKey, apiSecret string) string {
	s := fmt.Sprintf("%s:%s", apiKey, apiSecret)
	s = base64.StdEncoding.EncodeToString([]byte(s))

	return fmt.Sprintf("Basic %s", s)
}
