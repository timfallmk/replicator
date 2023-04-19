package cmd

// import (
// 	"net/url"

// 	"github.com/pterm/pterm"
// 	"github.com/sirupsen/logrus"
// 	"github.com/spf13/cobra"

// 	"ad.astra.com/gitlab/tss/replicator/pkg/provisioner"
// )

// var newCmd = &cobra.Command{
// 	Use:   "new",
// 	Short: "Provision a new stand",
// 	Long:  `Provision a new stand. This will create a new stand and configure it for use.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		pURL, err := url.Parse(provisionerURL)
// 		if err != nil {
// 			logrus.Fatal(err)
// 		}
// 		client, err := provisioner.New(*pURL, provisionerToken)
// 		if err != nil {
// 			logrus.Fatal(err)
// 		}
// 	},
// }
