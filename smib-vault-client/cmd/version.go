/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of application",
	Long: `Show version of application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Application version: 1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
