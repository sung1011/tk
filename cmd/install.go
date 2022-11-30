/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
		// mkdir
		// clone
		// zsh
		// alias
		// brew tap,install,cask
		// path
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
