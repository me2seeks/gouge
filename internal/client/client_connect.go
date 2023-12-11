package client

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/me2seeks/gouge/share/cnet"
	"github.com/me2seeks/gouge/share/settings"
	"golang.org/x/crypto/ssh"
)

func (c *Client) connect(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return errors.New("context canceled")
	default:
	}
	d := websocket.Dialer{
		HandshakeTimeout: time.Second * 5,
		TLSClientConfig:  c.tlsConfig,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		NetDialContext:   c.config.DialContext,
	}
	if c.proxyURL != nil {
		if err := c.setProxy(c.proxyURL, &d); err != nil {
			return err
		}
	}
	wsConn, _, err := d.DialContext(ctx, c.server, c.config.Headers)
	if err != nil {
		return err
	}
	conn := cnet.NewWebSocketConn(wsConn)
	c.Debugf("Handshaking...")
	sshConn, _, _, err := ssh.NewClientConn(conn, c.server, c.sshConfig)
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "unable to authenticate") {
			c.Infof("Authentication failed")
			c.Debugf(e)
		} else {
			c.Infof(e)
		}
		return err
	}
	defer sshConn.Close()
	c.Debugf("Sending config")
	sshConn.SendRequest("config", true, settings.EncodeConfig(c.computed))
	if err != nil {
		c.Infof("Config verification failed")
		return err
	}

	//TODO: bindSSH

	return nil
}
