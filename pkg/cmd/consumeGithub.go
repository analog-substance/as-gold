package cmd

import (
	"github.com/analog-substance/as-gold/pkg/gold"
	"github.com/analog-substance/as-gold/pkg/util"
	"github.com/analog-substance/fileutil"
	"github.com/spf13/cobra"
	"log"
)

var consumeGithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Checkout non-forked repositories for users/orgs",
	Long:  `This will checkout and process all non-forked repos`,
	Run: func(cmd *cobra.Command, args []string) {
		update, _ := cmd.Flags().GetBool("update")
		authToken, _ := cmd.Flags().GetString("auth-token")
		//repoJSONFile, _ := cmd.Flags().GetString("repos-json")

		solidGold.SetGitHubAccessToken(authToken)

		noOrgsMembers, _ := cmd.Flags().GetBool("no-org-members")
		noUserOrgs, _ := cmd.Flags().GetBool("no-user-orgs")

		if update {
			solidGold.UpdateGithub(!noOrgsMembers, !noUserOrgs)
		} else {
			organizations, _ := cmd.Flags().GetStringSlice("orgs")
			users, _ := cmd.Flags().GetStringSlice("users")
			orgFile, _ := cmd.Flags().GetString("org-file")
			userFile, _ := cmd.Flags().GetString("user-file")

			organizations = mergeFileWithSlice(orgFile, organizations)
			users = mergeFileWithSlice(userFile, users)

			solidGold.ConsumeGithubOrgs(!noOrgsMembers, organizations...)
			solidGold.ConsumeGithubUsers(!noUserOrgs, users...)
		}

		solidGold.ProcessPath(gold.GitHubFolderPath)

		if goldFile != "" {
			err := solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
			err = solidGold.ToJSONFile(goldFile)
			util.CheckErr(err)
		}
	},
}

func init() {
	consumeCmd.AddCommand(consumeGithubCmd)

	consumeGithubCmd.Flags().StringSliceP("orgs", "o", []string{}, "orgs(s) to search for")
	consumeGithubCmd.Flags().String("org-file", "", "file with organizations")
	consumeGithubCmd.Flags().String("user-file", "", "file with users")
	consumeGithubCmd.Flags().StringSliceP("users", "u", []string{}, "users(s) to search for")
	consumeGithubCmd.Flags().StringP("auth-token", "a", "", "github personal auth token")
	consumeGithubCmd.Flags().BoolP("no-org-members", "M", false, "When querying organizations, do not look at it's members")
	consumeGithubCmd.Flags().BoolP("no-user-orgs", "O", false, "When querying users, do not look at organizations they are members of")
	consumeGithubCmd.Flags().BoolP("update", "U", false, "Update local copies")
	//consumeGithubCmd.Flags().StringP("repos-json", "r", "", "repos json api response file")

}

func mergeFileWithSlice(file string, slice []string) []string {

	if file != "" {
		fileLines, err := fileutil.ReadLines(file)
		if err != nil {
			log.Printf("Error reading org-file %s: %s\n", file, err)
		} else {
			slice = append(slice, fileLines...)
		}
	}
	return slice
}
