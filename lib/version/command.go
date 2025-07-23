package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
