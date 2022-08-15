package cmd

import (
	"fmt"
	"github.com/analog-substance/gold/internal/gold"
	"github.com/analog-substance/gold/internal/util"

	"github.com/spf13/cobra"
)

var dedupeCmd = &cobra.Command{
	Use:   "dedupe",
	Short: "Deduplicate entries",
	Long: `This will look for entries with duplicate emails and names (if the
name contains a space) and merge them into a single record.`,
	Run: func(cmd *cobra.Command, args []string) {

		var solidGold *gold.SolidGold

		if goldFile != "" {
			solidGold = gold.FromJSONFile(goldFile)
		} else {
			solidGold = gold.NewSolidGold()
		}

		beforeMembers := len(solidGold.Members)
		solidGold.Group.MergeDuplicate()
		afterMembers := len(solidGold.Members)

		if goldFile != "" {
			err := solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
		}
		fmt.Printf("Before: %d\nAfter: %d\n", beforeMembers, afterMembers)
	},
}

func init() {
	rootCmd.AddCommand(dedupeCmd)
}
