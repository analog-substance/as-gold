package cmd

import (
	"github.com/analog-substance/as-gold/pkg/gold"
	"github.com/analog-substance/as-gold/pkg/util"
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
		noOrgsMembers, _ := cmd.Flags().GetBool("no-orgs-members")
		noUserOrgs, _ := cmd.Flags().GetBool("no-user-orgs")

		authToken, _ := cmd.Flags().GetString("auth-token")
		//repoJSONFile, _ := cmd.Flags().GetString("repos-json")

		solidGold.ConsumeGithubOrgs(!noOrgsMembers, authToken, orgs...)
		solidGold.ConsumeGithubUsers(!noUserOrgs, authToken, users...)

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
	consumeGithubCmd.Flags().StringP("auth-token", "a", "", "github personal auth token")
	consumeGithubCmd.Flags().BoolP("no-org-members", "M", false, "When querying organizations, do not look at it's members")
	consumeGithubCmd.Flags().BoolP("no-user-orgs", "O", false, "When querying users, do not look at organizations they are members of")
	//consumeGithubCmd.Flags().StringP("repos-json", "r", "", "repos json api response file")

}
