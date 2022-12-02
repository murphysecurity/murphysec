package api

import (
	"context"
	"crypto/tls"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/httplogger"
	"github.com/murphysecurity/murphysec/utils/must"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	ServerURL          string
	Logger             *zap.Logger
	EnableNetworkDebug bool
	Token              string
	AllowInsecure      bool
	Ctx                context.Context
}

func (c *Config) Build() (*Client, error) {
	// copied from http package
	var defaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     false, // disable force http2 attempt
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Allow insecure
	if c.AllowInsecure {
		defaultTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	var client Client
	client.ctx = c.Ctx

	client.client = &http.Client{Transport: defaultTransport}

	// API base URL
	if c.ServerURL == "" {
		client.baseUrl = must.A(url.Parse(env.DefaultServerURL))
	} else {
		u, e := url.Parse(c.ServerURL)
		if e != nil {
			return nil, e
		}
		client.baseUrl = u
	}

	// Logging & Debugging
	client.logger = zap.NewNop()
	if c.Logger != nil {
		client.logger = c.Logger
	}
	if c.EnableNetworkDebug {
		client.client.Transport = httplogger.New(client.client.Transport, c.Logger.Named("net"))
	}

	// Token
	client.token = c.Token

	return &client, nil
}
