package monica

import (
	"context"

	"github.com/oklog/run"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

type Name string

type Module interface{}

type Nameable interface {
	Name() Name
}

type ProviderCommand interface {
	ProvideCommand() *cobra.Command
}

type ProviderRunGroup interface {
	ProvideRunGroup(g *run.Group)
}

type ProviderCronJob interface {
	ProvideCronJob(c *cron.Cron)
}

type ProviderShutdown interface {
	ProvideShutdown() func(ctx context.Context) error
}

type Runnable interface {
	Run(ctx context.Context) error
}
