package cmd

import (
	"github.com/analog-substance/as-gold/pkg/gold"

	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge two or more entries by UUID",
	Long: `Merge two or more entries by UUID.

gold merge ff3ce69c-c054-473b-bfa6-9f0510383969 1bdd0e74-e78e-4008-80e2-31b6ca1c8352 053efb37-0349-4562-9409-680b15c65065`,
	Run: func(cmd *cobra.Command, args []string) {

		var solidGold *gold.SolidGold

		if goldFile != "" {
			solidGold = gold.FromJSONFile(goldFile)
		} else {
			solidGold = gold.NewSolidGold()
		}

		if len(args) >= 2 {
			primaryHuman := args[0]
			toMerge := args[1:]

			solidGold.Merge(primaryHuman, toMerge...)

			if goldFile != "" {
				solidGold.ToJSONFile(goldFile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}
