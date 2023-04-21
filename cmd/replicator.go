package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "replicator",
	Short: "A tool for provisioning new stands",
	Long: `Replicator is a tool to automate the provisiong process for new stands.
It can be used in an interactive mode to manually execute each step in the process, or
options can be passed directly to create a new stand in a hands-off manner.`,
}

// var logLevel string
var provisionerURL string
var provisionerToken string

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&provisionerURL, "provision-server", "l", "", "The URL of the Provision (MAAS) server")
	rootCmd.PersistentFlags().StringVarP(&provisionerToken, "api-key", "t", "", "The API key for the Provision (MAAS) server")

	rootCmd.MarkPersistentFlagRequired("provision-server")
	rootCmd.MarkPersistentFlagRequired("api-key")

	// Flags for the "new" command
	newCmd.Flags().StringVarP(&standName, "name", "n", "", "Name of the stand to provision")
	newCmd.Flags().StringVar(&standFQDN, "fqdn", "", "Fully Qualified Domain Name of the stand to provision")
	newCmd.Flags().StringVarP(&standUserData, "cloud-config", "f", "", "Path to the cloud-config (user-data) file to use for the stand")

	// TODO: Allow providing either the name or the FQDN, but not both
	newCmd.MarkFlagRequired("name")
	newCmd.MarkFlagRequired("cloud-config")

	// rootCmd.AddCommand(createCmd)
	// 	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(eraseCmd)
}
