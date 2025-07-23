package monica

import (
	"github.com/spf13/cobra"
)

// Serve is a module that provides a command to serve the application
// it will start the http server and the grpc server
type Serve struct {
	core *Core
}

func (s *Serve) Name() Name {
	return "serve"
}

func (s *Serve) ProvideCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "serve the application",
		Run: func(cmd *cobra.Command, args []string) {
			s.core.c.Start()
			s.core.g.Run()
		},
	}
}
