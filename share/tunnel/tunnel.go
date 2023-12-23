package tunnel

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/armon/go-socks5"
	"github.com/me2seeks/gouge/share/cio"
	"github.com/me2seeks/gouge/share/cnet"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	cio.Logger
	inbound   string
	outbound  string
	socket    bool
	keeoplive bool
}

type Tunnel struct {
	Config
	activeMut      sync.Mutex
	activatingConn waitgroup
	activeConn     ssh.Conn
	connStats      cnet.ConnCount
	socksServer    *socks5.Server
}

func New(c Config) *Tunnel {
	c.Logger.Fork("tun")
	t := &Tunnel{
		Config: c,
	}
	t.activatingConn.Add(1)
	extra := ""
	if t.socket {
		sl := log.New(io.Discard, "", 0)
		if t.Logger.Debug {
			sl = log.New(os.Stdout, "[socks]", log.Ldate|log.Ltime)
		}
		t.socksServer, _ = socks5.New(&socks5.Config{Logger: sl})
		extra += "(SOCKS enabled)"
	}
	t.Debugf("Created%s", extra)
	return t
}
