/*
Copyright © 2022 NAME HERE 

*/
package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"smib-vault-client/pkg/vault"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "List secret of backend.",
	Long: `Read command allow to get list of all secrets located on vault "path".
You could store received secrets to file. For this, please use flag "filename".
For the security of viewing information, there is opportunity to convert secret to encrypted by using AES, as a encryption key use global variables ENCRYPT_KEY. For this, please use flag "encrypt".
Or with flag "base64" data will be encode to base64.
Use only one flag: "encrypt" or "base64".
Flag separator let use different symbols which will be displaed between keys and values. By default separator is ': '"
Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 
"{VAULT_PROJECT_NAME}/{path}".
To the run command:
$ vault-client read --path vault-client/dso --filename SuperPuper.secret --base64 --separator "="`,

	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		filename, _ := cmd.Flags().GetString("filename")
		encrypt, _ := cmd.Flags().GetBool("encrypt")
		base64, _ := cmd.Flags().GetBool("base64")
		separator, _ := cmd.Flags().GetString("separator")
		if (base64 && encrypt) { 
			fmt.Println("Error. Can use only one flag: encrypt or base64.")
			os.Exit(1)
		}
		vault.Read(path, filename, encrypt, base64, separator)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.PersistentFlags().String("path", "", "Vault path where stored secrets")
	readCmd.MarkPersistentFlagRequired("path")
	readCmd.PersistentFlags().String("filename", "", "File to store secrets")
	readCmd.Flags().BoolP("encrypt", "", false, "Encrypt secret to aes256")
	readCmd.Flags().BoolP("base64", "", false, "Encrypt secret value to base64")
	readCmd.Flags().String("separator", ": ", "Symbols which will be serapate keys and values. By default separator is ': '")
}
