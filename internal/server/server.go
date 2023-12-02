package server

import (
	"context"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/me2seeks/gouge/share/ccrypto"
	"github.com/me2seeks/gouge/share/cio"
	"github.com/me2seeks/gouge/share/cnet"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	KeySeed   string
	KeyFile   string
	Proxy     string
	KeepAlive bool
}

type Server struct {
	*cio.Logger
	config      *Config
	fingerprint string
	httpServer  *cnet.HTTPServer
	sshConfig   *ssh.ServerConfig
	sessCount   int32
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  0,
	WriteBufferSize: 0,
}

func NewServer(Config *Config) (*Server, error) {
	server := &Server{
		Logger:     cio.NewLogger("server"),
		config:     Config,
		httpServer: cnet.NewHTTPServer(),
	}

	server.Info = true
	server.Debug = true

	var pemBytes []byte
	var err error
	pemBytes, err = ccrypto.GeneratePEM()
	if err != nil {
		log.Fatal(err)
	}
	private, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatal("Failed to parse key")
	}
	server.fingerprint = ccrypto.FingerPrint(private.PublicKey())
	server.sshConfig = &ssh.ServerConfig{
		ServerVersion:    "SSH-Gouge",
		PasswordCallback: server.authUser,
	}
	server.sshConfig.AddHostKey(private)

	return server, nil

}
func (s *Server) StartContext(ctx context.Context, host, port string) error {

	s.Infof("Fingerprint %s ", s.fingerprint)
	l, err := s.Listener(host, port)
	if err != nil {
		return err
	}
	h := http.Handler(http.HandlerFunc(s.handleClientHandler))

	return s.httpServer.GoServe(ctx, l, h)
}

func (s *Server) Wait() error {
	return s.httpServer.Wait()
}

func (s *Server) authUser(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	return nil, nil
}
