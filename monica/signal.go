package monica

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Signal struct {
	core *Core
}

func (s *Signal) Name() Name {
	return "signal"
}

func (s *Signal) Run(ctx context.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	log.Println("signal received, shutting down")
	s.core.Shutdown(ctx)
	return nil
}
