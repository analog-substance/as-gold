package cmd

import (
	"github.com/analog-substance/gold/internal/gold"
	"github.com/analog-substance/gold/internal/util"
	"github.com/spf13/cobra"
)

var consumeGithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Checkout non-forked repositories for users/orgs",
	Long:  `This will checkout and process all non-forked repos`,
	Run: func(cmd *cobra.Command, args []string) {

		var solidGold *gold.SolidGold

		if goldFile != "" {
			solidGold = gold.FromJSONFile(goldFile)
		} else {
			solidGold = gold.NewSolidGold()
		}

		orgs, _ := cmd.Flags().GetStringSlice("orgs")
		users, _ := cmd.Flags().GetStringSlice("users")

		solidGold.ConsumeGithubOrgs(true, orgs...)
		solidGold.ConsumeGithubUsers(true, users...)

		solidGold.ProcessPath("github.com")

		if goldFile != "" {
			solidGold.ToJSONFile(goldFile)
			err := solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
		}
	},
}

func init() {
	consumeCmd.AddCommand(consumeGithubCmd)

	consumeGithubCmd.Flags().StringSliceP("orgs", "o", []string{}, "orgs(s) to search for")
	consumeGithubCmd.Flags().StringSliceP("users", "u", []string{}, "users(s) to search for")

}
