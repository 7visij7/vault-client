/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"smib-vault-client/pkg/vault"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy all secrets from source to destinaton KV.",
	Long: `Copy command allow to dublicate secrets located in "source" to other KV which you have to define as "destination".
	If flag "soft" setted, and secret destination has already exist, coping does not run.
	Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 
		"{VAULT_PROJECT_NAME}/{path}".
	To the run command: 
	$ vault-client copy  --source smib/vault-client --destination smbu/vault-client`,
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")
		destination, _ := cmd.Flags().GetString("destination")
		soft, _ := cmd.Flags().GetBool("soft")
		vault.Copy(source, destination, soft)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
	copyCmd.PersistentFlags().String("source", "", "Vault path where stored secrets which need to copy")
	copyCmd.MarkPersistentFlagRequired("source")
	copyCmd.PersistentFlags().String("destination", "", "Vault path where need to put secrets")
	copyCmd.MarkPersistentFlagRequired("destination")
	copyCmd.Flags().BoolP("soft", "", false, "If secret exists, no replacement is performed")
}
