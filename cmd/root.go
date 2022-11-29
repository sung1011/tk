package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sung1011/tk/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tk",
	Short: "tickles cmd tools",
	Long:  `tk is a CLI for GO that tickles applications`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if len(os.Args) > 1 {
			path, errExec := exec.LookPath(os.Args[1])
			if errExec == nil {
				cmd := exec.Command(path)
				utils.RunCmd(cmd)
			}
		}

	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tk.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const CONFIG_FILE_DEFAULT = ".tk.yaml"

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// args
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		curPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// local
		pathExists, err := utils.PathExists(fmt.Sprintf("%s/%s", curPath, CONFIG_FILE_DEFAULT))
		if err != nil {
			panic(err)
		}
		if pathExists {
			viper.AddConfigPath(curPath)
			viper.SetConfigName(".tk")
		} else {
			// home
			pathExists, err = utils.PathExists(fmt.Sprintf("%s/%s", home, CONFIG_FILE_DEFAULT))
			if err != nil {
				panic(err)
			}
			if !pathExists {
				if err != nil {
					panic(err)
				}
				utils.CopyFile(fmt.Sprintf("%s/%s.default", curPath, CONFIG_FILE_DEFAULT), fmt.Sprintf("%s/%s", home, CONFIG_FILE_DEFAULT))
			}
			viper.AddConfigPath(home)
			viper.SetConfigName(".tk")
		}
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		panic(err)
	}
}
