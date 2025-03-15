package pool

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

// GetHttpClient returns a new http.Client that is configured to work well in
// a high concurrency environment. It will use the proxy from the environment
// and dial a connection with a 30 second timeout, and a keepalive of 30
// seconds. It will also keep up to 200 idle connections open to the same host,
// and to any host. It will also timeout any idle connections after 90
// seconds. TLS handshakes have a 10 second timeout and expect continue
// responses have a 1 second timeout.
func GetHttpClient() *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   200,
		MaxIdleConns:          200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{Transport: tr}

}

// StressEndpoint sends a request to the given url with the given method and payload
// and checks the status code of the response.
//
// It returns an error if the request fails or the status code is not 200.
func StressEndpoint(method string, url string, payload string) error {
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		slog.Error("TestEndpoint", "msg", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := GetHttpClient().Do(req)
	if err != nil {
		slog.Error("TestEndpoint Do", "msg", err.Error())
		return err
	}
	if res.StatusCode != 200 {
		slog.Error("TestEndpoint StatusCode", "msg", res.StatusCode)
		return fmt.Errorf("TestEndpoint StatusCode: %d", res.StatusCode)
	}
	defer res.Body.Close()
	io.Copy(io.Discard, res.Body)
	return nil
}
