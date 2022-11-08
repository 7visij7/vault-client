/*
Copyright © 2022 NAME HERE 

*/
package cmd

import (
	"github.com/spf13/cobra"
	"smib-vault-client/pkg/vault"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all secret of KV.",
	Long: `Delete command allow to remove all secrets located on vault "path".
	Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 
		"{VAULT_PROJECT_NAME}/{path}".
	To the run command: 
	$ vault-client delete --path smib/smib-vault-client`,
	
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		vault.Delete(path)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().String("path", "", "Vault path where stored secrets")
	deleteCmd.MarkPersistentFlagRequired("path")
}
