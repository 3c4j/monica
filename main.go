package main

import (
	"context"
	"fmt"
	"log"

	"github.com/3c4j/monica/lib/version"
	"github.com/3c4j/monica/monica"
	"github.com/3c4j/monica/pkg/logger"
	"github.com/3c4j/monica/user"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	GitHash           string
	GitBranch         string
	BuildTime         string
	LastCommitHash    string
	LastCommitBranch  string
	LastCommitTime    string
	LastCommitAuthor  string
	LastCommitMessage string
)

var (
	configPath = ""
	rootCmd    = &cobra.Command{
		Use:     "monica",
		Short:   "monica is a simple and easy-to-use monorepo manager",
		Version: fmt.Sprintf("%s@%s", GitHash, GitBranch),
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./config.yaml", "config file (default is ./config.yaml)")
}

func main() {
	err := monica.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	logger := logger.NewLogger(viper.GetString("log.level"), viper.GetString("log.format"), viper.GetString("log.output"))
	core := monica.NewCore(rootCmd, logger)
	core.Register(&version.Version{
		GitHash:   GitHash,
		GitBranch: GitBranch,
		BuildTime: BuildTime,
		Commit: version.GitCommitInfo{
			Hash:    LastCommitHash,
			Branch:  LastCommitBranch,
			Time:    LastCommitTime,
			Author:  LastCommitAuthor,
			Message: LastCommitMessage,
		},
	})
	core.RegisterWithFunc(user.InitModule)
	if err := core.Run(context.Background()); err != nil {
		log.Fatalf("Fatal error core: %s \n", err)
	}
}
