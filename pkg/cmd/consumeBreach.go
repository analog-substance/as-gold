package cmd

import (
	"github.com/analog-substance/as-gold/pkg/util"
	"github.com/spf13/cobra"
)

var consumeBreachCmd = &cobra.Command{
	Use:   "breach",
	Short: "Consume breach dump data",
	Long: `Consume breach dump data from a TSV file. For example:

username@place.tld	password!
otheruser@place.tld	Spring 2022!`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			solidGold.ConsumeBreachFiles(args...)
		}
		if goldFile != "" {
			err := solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
		}
	},
}

func init() {
	consumeCmd.AddCommand(consumeBreachCmd)
}
