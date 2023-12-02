package cos

import (
	"context"
	"os"
	"os/signal"
)

func InterruptContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		signal.Stop(sig)
		cancel()
	}()
	return ctx
}
