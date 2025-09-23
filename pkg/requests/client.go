package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	urllib "net/url"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type HttpClient struct {
	log    *logrus.Entry
	tracer trace.Tracer
}

func NewHttpClient(logger *logrus.Entry) *HttpClient {
	return &HttpClient{
		log:    logger.Dup().WithField("package", packageName),
		tracer: otel.Tracer(packageName),
	}
}

func (c *HttpClient) sendTrace(t time.Time, url string, remoteAddr string, method string, status int, reqId string, e error) {
	h := "-"
	path := "-"
	q := "-"
	u, err := urllib.Parse(url)
	if err == nil {
		h = u.Host
		path = u.Path
		if u.RawQuery != "" {
			q = u.RawQuery
		}
	}
	lvl := logrus.InfoLevel
	if e != nil {
		lvl = logrus.ErrorLevel
	}
	s, _ := os.Hostname()
	if s == "" {
		s = "-"
	}
	if reqId == "" {
		reqId = "-"
	}
	c.log.Dup().WithTime(t).WithField("error", e).Logf(
		lvl,
		`server_name=%s request_id=%s remote_host=%s remote_addr=%s method=%s path=%s query=%s status=%d latenct_ms=%d`,
		s,
		reqId,
		h,
		remoteAddr,
		method,
		path,
		q,
		status,
		time.Since(t).Milliseconds(),
	)
}

func (c *HttpClient) Request(ctx context.Context, method string, url string, headers map[string]string, body io.Reader, timeout int) (*Response, error) {
	return c.RequestWithTLSConfig(ctx, method, url, headers, body, timeout, nil)
}

func (c *HttpClient) RequestWithTLSConfig(ctx context.Context, method string, url string, headers map[string]string, body io.Reader, timeout int, tlsCfg *tls.Config) (*Response, error) {
	ctx, span := c.tracer.Start(ctx, "requests.RequestWithTLSConfig")
	defer span.End()

	timeBegin := time.Now()
	if timeout == 0 {
		timeout = Timeout
	}
	r := Response{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		c.sendTrace(timeBegin, url, "-", method, 0, util.GetHttpRequestId(ctx), err)
		return &r, err
	}
	var remoteAddr string
	ct := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr().String()
		},
	}
	req = req.WithContext(ctx)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), ct))

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if tlsCfg == nil {
		tlsCfg = &DefaultTLSConfig
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(
			&http.Transport{
				TLSClientConfig: tlsCfg,
			},
			otelhttp.WithTracerProvider(otel.GetTracerProvider()),
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		),
		Timeout: time.Duration(timeout) * time.Second,
	}

	// TODO: Retry should be here
	// Example Cr. https://stackoverflow.com/questions/50676817/does-the-http-request-automatically-retry
	resp, err := client.Do(req)
	if err != nil {
		c.sendTrace(timeBegin, url, remoteAddr, method, 0, util.GetHttpRequestId(ctx), err)
		return &r, err
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)

	r.Code = resp.StatusCode
	r.Status = resp.Status
	r.Body = buf.Bytes()
	r.Header = resp.Header
	_ = resp.Body.Close()

	c.sendTrace(timeBegin, url, remoteAddr, method, r.Code, util.GetHttpRequestId(ctx), err)
	return &r, nil
}

func (c *HttpClient) Get(ctx context.Context, url string, headers map[string]string, body io.Reader, timeout int) (*Response, error) {
	return c.Request(ctx, http.MethodGet, url, headers, body, timeout)
}

func (c *HttpClient) Post(ctx context.Context, url string, headers map[string]string, body io.Reader, timeout int) (*Response, error) {
	return c.Request(ctx, http.MethodPost, url, headers, body, timeout)
}

func (c *HttpClient) Put(ctx context.Context, url string, headers map[string]string, body io.Reader, timeout int) (*Response, error) {
	return c.Request(ctx, http.MethodPut, url, headers, body, timeout)
}

func (c *HttpClient) Delete(ctx context.Context, url string, headers map[string]string, body io.Reader, timeout int) (*Response, error) {
	return c.Request(ctx, http.MethodDelete, url, headers, body, timeout)
}
