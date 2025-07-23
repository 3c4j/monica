package main

import (
	"context"
	"fmt"
	"log"

	"github.com/3c4j/monica/monica"
	"github.com/spf13/cobra"
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

	core := monica.NewCore(rootCmd)
	core.RegisterModule(&Version{
		GitHash:   GitHash,
		GitBranch: GitBranch,
		BuildTime: BuildTime,
		Commit: GitCommitInfo{
			Hash:    LastCommitHash,
			Branch:  LastCommitBranch,
			Time:    LastCommitTime,
			Author:  LastCommitAuthor,
			Message: LastCommitMessage,
		},
	})
	if err := core.Run(context.Background()); err != nil {
		log.Fatalf("Fatal error core: %s \n", err)
	}
}

type Version struct {
	GitHash   string
	GitBranch string
	BuildTime string
	Commit    GitCommitInfo
}

type GitCommitInfo struct {
	Hash    string
	Branch  string
	Time    string
	Author  string
	Message string
}

func (v *Version) ProvideCommand() *cobra.Command {
	detail := false
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version of the application",
		Run: func(cmd *cobra.Command, args []string) {
			if v.GitHash == "" {
				v.GitHash = "unknown"
			}
			if v.GitBranch == "" {
				v.GitBranch = "unknown"
			}
			if v.BuildTime == "" {
				v.BuildTime = "unknown"
			}
			if detail {
				fmt.Println("Git hash:", v.GitHash)
				fmt.Println("Git branch:", v.GitBranch)
				fmt.Println("Build at:", v.BuildTime)
				fmt.Println("Last Commit:")
				fmt.Println("\tHash:", v.Commit.Hash)
				fmt.Println("\tBranch:", v.Commit.Branch)
				fmt.Println("\tTime:", v.Commit.Time)
				fmt.Println("\tAuthor:", v.Commit.Author)
				fmt.Println("\tMessage:", v.Commit.Message)
			} else {
				fmt.Println("Version:", v.GitHash+"@"+v.GitBranch)
				fmt.Println("Build at:", v.BuildTime)
			}
		},
	}
	cmd.Flags().BoolVarP(&detail, "detail", "d", false, "print the detail version of the application")
	return cmd
}
