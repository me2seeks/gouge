package cnet

import (
	"context"
	"errors"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type HTTPServer struct {
	*http.Server
	waiter *errgroup.Group
	//listenErr error
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		Server: &http.Server{},
	}
}

func (s *HTTPServer) ListenAndServe() error {
	return nil
}

func (h *HTTPServer) GoServe(ctx context.Context, l net.Listener, handler http.Handler) error {
	if ctx == nil {
		panic("nil context")
	}
	h.Handler = handler
	h.waiter, ctx = errgroup.WithContext(ctx)
	h.waiter.Go(func() error {
		return h.Serve(l)
	})
	go func() {
		<-ctx.Done()
		h.Close()
	}()
	return nil
}
func (h *HTTPServer) Close() error {
	return nil
}

func (h *HTTPServer) Wait() error {
	if h.waiter == nil {
		return errors.New("nil ctx")
	}
	err := h.waiter.Wait()
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}
