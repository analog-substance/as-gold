package cmd

import (
	"github.com/analog-substance/as-gold/pkg/util"
	"github.com/spf13/cobra"
)

var consumeGophishCmd = &cobra.Command{
	Use:   "gophish",
	Short: "Consume Gophish CSV",
	Long:  `Consume Gophish CSV`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			solidGold.ConsumeGophishFiles(args...)
		}
		if goldFile != "" {
			err := solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
		}
	},
}

func init() {
	consumeCmd.AddCommand(consumeGophishCmd)
}
