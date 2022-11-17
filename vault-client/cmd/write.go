/*
Copyright © 2022 NAME HERE 

*/
package cmd

import (
	"github.com/spf13/cobra"
	"smib-vault-client/pkg/vault"
	"smib-vault-client/pkg/common"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write secrets to the backend.",
	Long: `To use "write" command you should prepare file with secret like this example.yaml:
smib: "foo"
smbn: "bar"
smbu: "secret"

Next parametr which you shouls specify is vault "path" where stored secrets. Exp: "smib/smib-vault-client". 
With flag "replace" application changes value to vault path where secrets stored.
Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 
	"{VAULT_PROJECT_NAME}/{path}".
To the run command: 
$ vault-client write --fileaname example.yaml --path smib-vault-client/dso --replace`,

	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		filename, _ := cmd.Flags().GetString("filename")
		replace, _ := cmd.Flags().GetBool("replace")
		secrets := common.GetSecrets(filename)
		vault.Write(path, secrets)
		if replace {common.AddAnchorToFile(path, filename, secrets)}
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
	writeCmd.PersistentFlags().String("path", "", "Vault path where stored secrets")
	writeCmd.MarkPersistentFlagRequired("path")
	writeCmd.PersistentFlags().String("filename", "", "File where stored secrets")
	writeCmd.MarkPersistentFlagRequired("filename")
	writeCmd.Flags().BoolP("replace", "", false, "Flag to replace value to vault path where secrets stored.")
}
