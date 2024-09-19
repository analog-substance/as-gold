package cmd

import (
	"github.com/spf13/cobra"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume data from source via subcommand",
	Long: `Consume data from various sources. For example:

gold consume git path/to/gitrepo`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("dig called")
	//},
}

func init() {
	RootCmd.AddCommand(consumeCmd)
}
