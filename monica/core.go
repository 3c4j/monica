package monica

import (
	"context"
	"log"
	"reflect"

	"github.com/3c4j/monica/pkg/logger"
	"github.com/oklog/run"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

type Core struct {
	modules map[Name]Module
	cmd     *cobra.Command
	g       *run.Group
	c       *cron.Cron
	logger  *logger.Logger
}

func NewCore(rootCmd *cobra.Command, logger *logger.Logger) *Core {
	core := &Core{
		modules: make(map[Name]Module),
		cmd:     rootCmd,
		g:       &run.Group{},
		c:       cron.New(cron.WithSeconds()),
		logger:  logger,
	}
	// register default modules
	core.Register(&Serve{core: core})
	core.Register(&Signal{core: core})
	return core
}

func (c *Core) Register(m Module) {
	name := ""
	if nameable, ok := m.(Nameable); ok {
		name = string(nameable.Name())
	} else {
		// 使用反射获取模块名
		typ := reflect.TypeOf(m).Elem()
		name = typ.String()
	}
	if _, ok := c.modules[Name(name)]; ok {
		panic("module " + name + " already registered")
	}
	c.modules[Name(name)] = m
}

func (c *Core) RegisterWithFunc(fn func(*logger.Logger) (Module, error)) {
	m, err := fn(c.logger)
	if err != nil {
		panic(err)
	}
	c.Register(m)
}

func (c *Core) GetModule(name Name) Module {
	return c.modules[name]
}

func (c *Core) Run(ctx context.Context) error {
	for _, m := range c.modules {
		if mt, ok := m.(ProviderCommand); ok {
			c.cmd.AddCommand(mt.ProvideCommand())
		}

		if mt, ok := m.(ProviderRunGroup); ok {
			mt.ProvideRunGroup(c.g)
		}

		if mt, ok := m.(ProviderCronJob); ok {
			mt.ProvideCronJob(c.c)
		}

		if mt, ok := m.(Runnable); ok {
			_ctx, cancel := context.WithCancel(ctx)
			c.g.Add(func() error {
				return mt.Run(_ctx)
			}, func(err error) {
				cancel()
			})
		}
	}

	return c.cmd.Execute()
}

func (c *Core) Shutdown(ctx context.Context) error {
	for name, m := range c.modules {
		if mt, ok := m.(ProviderShutdown); ok {
			log.Println("shutting down module", name)
			if err := mt.Shutdown(ctx); err != nil {
				log.Println("error shutting down module", name, err)
			}
		}
	}
	c.c.Stop()
	return nil
}
