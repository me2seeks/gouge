package server

import (
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/me2seeks/gouge/share/cnet"
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

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddInt32(&s.sessCount, 1)
	l := s.Fork("session#%d", id)

	l.Infof("Websocket request %s", r.URL.Path)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.Debugf("Failed to upgrade websocket (%s)", err)
		return
	}

	conn := cnet.NewWebSocketConn(wsConn)
	// perform SSH handshake on net.Conn
	l.Debugf("Handshaking with %s...", r.RemoteAddr)
	sshConn, _, _, err := ssh.NewServerConn(conn, s.sshConfig)
	if err != nil {
		l.Debugf("Failed to handshake (%s)", err)
		return
	}

	// _, b, _ := sshConn.SendRequest("ping", true, nil)
	// println(string(b))
	time.Sleep(20 * time.Second)
	sshConn.Close()
}
