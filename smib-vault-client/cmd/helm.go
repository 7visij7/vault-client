/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"smib-vault-client/pkg/helm"
)

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Seek yaml files in directory and replace necessary secrets from Vault.",
	Long: `Helm command allow to find all yaml files in directoy, and replace values from Vault where you set anchor like this:
Variale_One: VAULT:smib/smib-common-secret#KAFKA_PASSWORD
Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server.
By default vault path will be formed: "{VAULT_PROJECT_NAME}/{ris}/{servicename}#{key}",
    ris -name of information system, servicename - name of service, key - name of variable.
Also you can specify directory where seek yaml files, for this use flag "path".
By default value of varaible will be encode to Base64, if you want to get value as is - use flag "raw". 

To the run command:
$ vault-client helm --path ./secret`,

	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		raw, _ := cmd.Flags().GetBool("raw")
        helm.Run(path, raw)
	},
}

func init() {
	rootCmd.AddCommand(helmCmd)
	helmCmd.PersistentFlags().String("path", "./", "Parameter path define catalog where app will search all yaml files. By default path set the current directory")
	helmCmd.Flags().BoolP("raw", "", false, "Not encrypt value to base64.")
}
