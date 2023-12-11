package server

import (
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/me2seeks/gouge/share/cnet"
	"github.com/me2seeks/gouge/share/settings"
	"golang.org/x/crypto/ssh"
)

func (s *Server) handleClientHandler(w http.ResponseWriter, r *http.Request) {

	upgrade := strings.ToLower(r.Header.Get("Upgrade"))
	if upgrade == "websocket" {
		s.handleWebsocket(w, r)
		return
	}

	w.WriteHeader(404)
	w.Write([]byte("Not found"))

}

func (s *Server) handleWebsocket(w http.ResponseWriter, req *http.Request) {
	id := atomic.AddInt32(&s.sessCount, 1)
	l := s.Fork("session#%d", id)

	l.Infof("Websocket request %s", req.URL.Path)
	wsConn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		l.Debugf("Failed to upgrade websocket (%s)", err)
		return
	}

	conn := cnet.NewWebSocketConn(wsConn)
	// perform SSH handshake on net.Conn
	l.Debugf("Handshaking with %s...", req.RemoteAddr)
	sshConn, _, reqs, err := ssh.NewServerConn(conn, s.sshConfig)
	if err != nil {
		l.Debugf("Failed to handshake (%s)", err)
		return
	}
	var r *ssh.Request
	select {
	case r = <-reqs:
	case <-time.After(10 * time.Second):
		l.Debugf("Timeout waiting for configuration")
		sshConn.Close()
		return
	}
	failed := func(err error) {
		l.Debugf("Failed: %s", err)
		r.Reply(false, []byte(err.Error()))
	}
	if r.Type != "config" {
		failed(s.Errorf("expecting config request"))
		return
	}
	c, err := settings.DecodeConfig(r.Payload)
	if err != nil {
		failed(s.Errorf("invalid config"))
		return
	}
	println(c.Remotes)
	r.Reply(true, nil)

	sshConn.Close()
}
