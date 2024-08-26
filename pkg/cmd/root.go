package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime/debug"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var goldFile string

func SetVersionInfo(versionStr, commitStr string) {
	buildInfo, ok := debug.ReadBuildInfo()
	buildType := "unknown"
	if ok {
		if versionStr != "v0.0.0" {
			// goreleaser must have set the version
			// lets add gh to the end so we know this release came from github
			buildType = "release"
		} else {
			// not a goreleaser build. lets grab build info from build settings
			versionStr = buildInfo.Main.Version

			if buildInfo.Main.Version == "(devel)" {
				for _, bv := range buildInfo.Settings {
					if bv.Key == "vcs.revision" {
						commitStr = bv.Value[0:8]
						buildType = "go-local"
						break
					}
				}
			} else {
				buildType = "go-remote"
				commitStr = buildInfo.Main.Version
			}
		}
	} else {
		log.Println("Version info not found in build info")
	}

	if os.Getenv("DEBUG_BUILD_INFO") == "1" {
		fmt.Println(buildInfo)
	}

	rootCmd.Version = fmt.Sprintf("%s (%s@%s)", versionStr, buildType, commitStr)
}

var rootCmd = &cobra.Command{
	Use:   "as-gold",
	Short: "Extract valuable human data from various sources",
	Long:  `Combine data from git repos, breach dumps into a collection of people.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gold.yaml)")
	rootCmd.PersistentFlags().StringVar(&goldFile, "gold", "solid-gold.json", "gold data file")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
