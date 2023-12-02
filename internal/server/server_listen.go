package server

import "net"

func (s *Server) Listener(host, port string) (net.Listener, error) {
	//TODO TLS
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}

	s.Infof("Listening on %s:%s", host, port)

	return l, nil
}
