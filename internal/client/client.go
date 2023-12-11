package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/me2seeks/gouge/share/ccrypto"
	"github.com/me2seeks/gouge/share/cio"
	"github.com/me2seeks/gouge/share/cnet"
	"github.com/me2seeks/gouge/share/settings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Fingerprint      string
	Server           string
	Proxy            string
	Auth             string
	Remotes          []string
	KeepAlive        time.Duration
	MaxRetryCount    int
	MaxRetryInterval time.Duration
	Headers          http.Header
	DialContext      func(ctx context.Context, network, addr string) (net.Conn, error)
	TLS              TLSConfig
	Verbose          bool
}

type TLSConfig struct {
	InsecureSkipVerify bool
	CA                 string
	Key                string
	Cert               string
	ServerName         string
}

type Client struct {
	*cio.Logger
	config    *Config
	computed  settings.Config
	sshConfig *ssh.ClientConfig
	tlsConfig *tls.Config
	proxyURL  *url.URL
	connCount cnet.ConnCount
	server    string
	eg        *errgroup.Group
	stop      func()
}

func NewClient(c *Config) (*Client, error) {
	if !strings.HasPrefix(c.Server, "http://") {
		c.Server = "http://" + c.Server
	}
	if c.MaxRetryInterval < time.Minute {
		c.MaxRetryInterval = time.Minute
	}
	u, err := url.Parse(c.Server)
	if err != nil {
		return nil, err
	}

	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)
	if !regexp.MustCompile(`:\d+$`).MatchString(u.Host) {
		if u.Scheme == "wss" {
			u.Host += ":443"
		} else {
			u.Host += ":80"
		}
	}

	hasReverse := false
	hasSocks := false
	hasStdio := false

	client := &Client{
		Logger:    cio.NewLogger("client"),
		config:    c,
		server:    u.String(),
		tlsConfig: nil,
	}
	client.Logger.Info = true

	if u.Scheme == "wss" {
		//TODO: TLSConfig
	}

	for _, s := range c.Remotes {
		r, err := settings.DecodeRemote(s)
		if err != nil {
			return nil, fmt.Errorf("failed to decode remote '%s': %s", s, err)
		}
		if r.Socks {
			hasSocks = true
		}
		if r.Reverse {
			hasReverse = true
		}
		if r.Stdio {
			if hasStdio {
				return nil, errors.New("only one stdio is allowed")
			}
			hasStdio = true
		}
		//confirm non-reverse tunnel is available
		if !r.Reverse && !r.Stdio && !r.CanListen() {
			return nil, fmt.Errorf("Client cannot listen on %s", r.String())
		}
		client.computed.Remotes = append(client.computed.Remotes, r)
	}

	println(hasReverse, hasSocks, hasStdio)

	if p := c.Proxy; p != "" {
		client.proxyURL, err = url.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL (%s)", err)
		}
	}
	user, pass := settings.ParseAuth(c.Auth)
	client.sshConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: client.verifyServer,
		Timeout:         30 * time.Second,
	}

	return client, nil
}

func (c *Client) StartContext(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.stop = cancel
	c.eg, ctx = errgroup.WithContext(ctx)
	via := ""
	if c.proxyURL != nil {
		via = " via " + c.proxyURL.String()
	}
	c.Infof("Connecting to %s via %s", c.server, via)
	//connent to server
	c.eg.Go(func() error {
		return c.connect(ctx)
	})

	return nil
}

func (c *Client) Wait() error {
	return c.eg.Wait()
}

func (c *Client) verifyServer(hostname string, remote net.Addr, key ssh.PublicKey) error {
	if c.config.Fingerprint != "" {
		if c.config.Fingerprint != ccrypto.FingerPrint(key) {
			return fmt.Errorf("fingerprint mismatch")
		}
		return nil
	}
	return nil
}

func (c *Client) setProxy(u *url.URL, d *websocket.Dialer) error {
	//TODO
	return nil
}
