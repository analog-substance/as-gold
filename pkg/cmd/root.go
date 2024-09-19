package cmd

import (
	"fmt"
	"github.com/analog-substance/as-gold/pkg/gold"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var goldFile string

var solidGold *gold.SolidGold

var RootCmd = &cobra.Command{
	Use:   "as-gold",
	Short: "Extract valuable human data from various sources",
	Long:  `Combine data from git repos, breach dumps into a collection of people.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if goldFile != "" {
			solidGold = gold.FromJSONFile(goldFile)
		} else {
			solidGold = gold.NewSolidGold()
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gold.yaml)")
	RootCmd.PersistentFlags().StringVar(&goldFile, "gold", "solid-gold.json", "gold data file")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".gold")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
