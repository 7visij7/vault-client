/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"smib-vault-client/pkg/discovery"
)

// discoveryCmd represents the discovery command
var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Seek yaml files in directory and put variables in section secret to Vault",
	Long: `Discovery command allow to find all yaml files in directoy "helm" and put secrets to Vault.
To use "discovery" you should specify flags "ris", "enviroment" and "servicename".
Flag "enviroment" must be one of [dso, st ... etc] values depends on in what catalog seek secrets.
With flag "replace" application changes value to vault path where secrets stored.
Have to difine enviroment variable $VAULT_PROJECT_NAME with name of KV on Vault server. Vault path will be: 
	"{VAULT_PROJECT_NAME}/{ris}/{servicename}#{key}".
To the run command:
	$ vault-client discovery --ris smib --servicename testing-new  --replace --enviroment dso`,

	Run: func(cmd *cobra.Command, args []string) {
		ris, _ := cmd.Flags().GetString("ris")
		servicename, _ := cmd.Flags().GetString("servicename")
		replace, _ := cmd.Flags().GetBool("replace")
		enviroment, _ := cmd.Flags().GetString("enviroment")

		discovery.Run(ris, servicename, replace, enviroment)
	},
}

func init() {
	rootCmd.AddCommand(discoveryCmd)
	discoveryCmd.PersistentFlags().String("ris", "", "Name of information system")
	discoveryCmd.MarkPersistentFlagRequired("ris")
	discoveryCmd.PersistentFlags().String("servicename", "", "Name of service")
	discoveryCmd.MarkPersistentFlagRequired("servicename")
	discoveryCmd.PersistentFlags().String("enviroment", "", "Enviroment must be one of [dso, st] values depends on where application running")
	discoveryCmd.MarkPersistentFlagRequired("enviroment")
	discoveryCmd.Flags().BoolP("replace", "", false, "Flag to replace value to vault path where secrets stored.")
}
