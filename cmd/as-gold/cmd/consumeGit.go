package cmd

import (
	"github.com/analog-substance/as-gold/internal/gold"
	"github.com/analog-substance/as-gold/internal/util"
	"github.com/spf13/cobra"
)

var consumeGitCmd = &cobra.Command{
	Use:   "git",
	Short: "Consume names and emails from Git repositories",
	Long:  `This will collect names and emails from author and committer information in git repositories.`,
	Run: func(cmd *cobra.Command, args []string) {

		var solidGold *gold.SolidGold

		if goldFile != "" {
			solidGold = gold.FromJSONFile(goldFile)
		} else {
			solidGold = gold.NewSolidGold()
		}

		if len(args) != 1 {
			panic("need a path")
		}

		solidGold.ProcessPath(args[0])

		if goldFile != "" {
			err := solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
		}
	},
}

func init() {
	consumeCmd.AddCommand(consumeGitCmd)
}
