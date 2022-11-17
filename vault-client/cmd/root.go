/*
Copyright © 2022 NAME HERE 

*/
package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vault-client",
	Short: "This is vault client of SMIB team.",
	Long: `This is client to work with Vault server.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}


