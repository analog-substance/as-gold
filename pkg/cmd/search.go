package cmd

import (
	"fmt"
	"github.com/analog-substance/as-gold/pkg/gold"
	"github.com/analog-substance/as-gold/pkg/util"
	"github.com/spf13/cobra"
)

// searchCmd represents the breach command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search entries by various values",
	Long: `Search entries by name, email, or if they have a password. Output JSON,
name, email, or credential dump. For Example:

Only return email addresses
	gold search -E

Return entries with example.tld email domains
	gold search -d example.tld

Return entries with "company" in them and save to new json file
	gold search -s company -S company-solid-gold.json

`,
	Run: func(cmd *cobra.Command, args []string) {

		printEmails, _ := cmd.Flags().GetBool("print-emails")
		printNames, _ := cmd.Flags().GetBool("print-names")
		printPasswords, _ := cmd.Flags().GetBool("print-passwords")
		printCredentials, _ := cmd.Flags().GetBool("print-credentials")

		saveFile, _ := cmd.Flags().GetString("save")
		domains, _ := cmd.Flags().GetStringSlice("domains")
		searchStrings, _ := cmd.Flags().GetStringSlice("search")
		onlyWithPasswords, _ := cmd.Flags().GetBool("passwords")

		if printCredentials {
			onlyWithPasswords = true
		}

		subGroup := solidGold.Group

		if len(domains) > 0 {
			subGroup = subGroup.FindWithEmailDomains(domains...)
		}

		if len(searchStrings) > 0 {
			subGroup = subGroup.FindWithString(searchStrings...)
		}

		if onlyWithPasswords {
			subGroup = subGroup.FindWithPasswords()
		}

		newSG := gold.NewSolidGold()
		newSG.Group = subGroup

		if saveFile != "" {
			err := newSG.ToJSONFile(saveFile)
			util.CheckErr(err)
		}

		if printEmails {
			for _, h := range newSG.Members {
				for _, e := range h.Emails {
					fmt.Println(e)
				}
			}
			return
		}
		if printNames {
			for _, h := range newSG.Members {
				for _, e := range h.Names {
					fmt.Println(e)
				}
			}
			return
		}
		if printPasswords {
			for _, h := range newSG.Members {
				for _, e := range h.Passwords {
					fmt.Println(e)
				}
			}
			return
		}
		if printCredentials {
			for _, h := range newSG.Members {
				for _, e := range h.Emails {
					for _, p := range h.Passwords {
						fmt.Printf("%s\t%s\n", e, p)
					}
				}
			}
			return
		}

		fmt.Println(string(newSG.ToJSON()))

	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringSliceP("domains", "d", []string{}, "domain(s) to search for")
	searchCmd.Flags().StringSliceP("search", "s", []string{}, "strings(s) to search for in name and email")
	searchCmd.Flags().StringP("save", "S", "", "save results to file")
	searchCmd.Flags().BoolP("passwords", "p", false, "return humans with passwords")

	searchCmd.Flags().BoolP("print-emails", "E", false, "only print emails of matched humans")
	searchCmd.Flags().BoolP("print-names", "N", false, "only print names of matched humans")
	searchCmd.Flags().BoolP("print-passwords", "P", false, "only print passwords of matched humans")
	searchCmd.Flags().BoolP("print-credentials", "C", false, "print credentials")

}
