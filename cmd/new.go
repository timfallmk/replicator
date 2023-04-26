package cmd

import (
	"net/url"
	"time"

	"github.com/maas/gomaasclient/entity"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"ad.astra.com/gitlab/tss/replicator/pkg/provisioner"
)

var standName string
var standFQDN string
var standUserData string

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Provision a new stand",
	Long:  `Provision a new stand. This will create a new stand and configure it for use.`,
	Run: func(cmd *cobra.Command, args []string) {
		pURL, err := url.Parse(provisionerURL)
		if err != nil {
			logrus.Fatal(err)
		}
		client, err := provisioner.New(*pURL, provisionerToken)
		if err != nil {
			logrus.Fatal(err)
		}

		// Construct FQDN if not provided
		if standFQDN == "" && standName != "" {
			// standFQDN = standName + ".ad.astra.com"
			standFQDN = standName + ".ad.astraspace.com"
		}

		progressBar, err := pterm.DefaultProgressbar.WithTotal(5).WithTitle("Provisioning new stand").Start()
		// newSpinner, err := pterm.DefaultSpinner.Start()
		if err != nil {
			logrus.Error(err)
		}
		progressBar.UpdateTitle("Starting deployment process...")
		// newSpinner.UpdateText("Starting deployment process...")

		// Deployment step through process
		// Get potential new candidates
		newList, err := client.FindNewMachines()
		if err != nil {
			// newSpinner.Fail("Failed to find any new machines")
			logrus.Fatal(err)
		}
		newListNames := make([]string, 0)
		for n := range newList {
			newListNames = append(newListNames, newList[n].Hostname)
		}
		if len(newListNames) == 0 {
			// newSpinner.Fail("No fresh machines found")
			logrus.Fatal("No fresh machines found")
		}

		// TODO: Check if already exists

		// newSpinner.Stop()
		pterm.Success.Println("Found fresh machines")
		progressBar.Increment()
		newSelectedMachine, err := pterm.DefaultInteractiveSelect.WithOptions(
			newListNames,
		).Show()
		if err != nil {
			// newSpinner.Fail("Failed to select a machine")
			logrus.Fatal(err)
		}
		pterm.Info.Printf("Selected machine %s\n", newSelectedMachine)
		// Reverse lookup
		// TODO: Make this better
		var selectedMachine entity.Machine
		for m := range newList {
			logrus.Debugf("Comparing %s to %s\n", newList[m].Hostname, newSelectedMachine)
			if newList[m].Hostname == newSelectedMachine {
				selectedMachine = newList[m]
				logrus.Debugf("Found matching ID %s for machine %+v \n", newList[m].SystemID, selectedMachine.FQDN)
				break
			}
		}
		if selectedMachine.SystemID == "" {
			logrus.Fatalf("Failed to reverse map machine %s to %+v \n", newSelectedMachine, "")
		}

		// DEBUG
		// newMachine := newList[0]

		progressBar.UpdateTitle("Comissioning machine...")
		// comissionSpinner, err := pterm.DefaultSpinner.Start("Comissioning machine...")
		if err != nil {
			logrus.Error(err)
		}
		// comissionSpinner.UpdateText("Setting hostname...\n")
		// Set the hostname
		logrus.Debugf("Attempting to set hostname to %s\n", standFQDN)
		err = client.SetHostname(selectedMachine.SystemID, standFQDN)
		if err != nil {
			// comissionSpinner.Fail("Failed to set hostname\n")
			logrus.Fatal(err)
		}
		// comissionSpinner.Success("Hostname set successfully\n")

		// comissionSpinner.UpdateText("Setting power type and comissioning...\n")
		// Comission
		// TODO: Get the updated stand information to make sure it matches what we expect
		node, comissionErr := client.CommissionNode(standFQDN)
		if comissionErr != nil {
			// comissionSpinner.Fail("Failed to comission node\n")
			logrus.Fatal(comissionErr)
		}

		// comissionSpinner.UpdateText("Waiting for comissioning to complete...\n")
		// comissionSpinner.Start()
		// Wait for comissioning to complete
		comissionWaitErr := client.WaitForComissioned(node)
		if comissionWaitErr != nil {
			// comissionSpinner.Fail("Failed to get a successful status from the comissioning process\n")
			logrus.Fatal(comissionWaitErr)
		}
		pterm.Success.Println("Comissioning completed successfully")
		progressBar.Increment()
		// comissionSpinner.Success("Comissioning completed successfully\n")
		// comissionSpinner.Stop()

		// Arbitrary wait until machine updates state
		// TODO: Remove this
		// waitSpinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Waiting for state to settle...")
		progressBar.UpdateTitle("Waiting for state to settle...")
		if err != nil {
			logrus.Error(err)
		}
		time.Sleep(20 * time.Second)
		pterm.Success.Println("State settled")
		progressBar.Increment()
		// waitSpinner.Success("Time elapsed\n")
		// waitSpinner.Stop()

		progressBar.UpdateTitle("Deploying node...")
		// deploySpinner, err := pterm.DefaultSpinner.Start("Comissioning machine...")
		if err != nil {
			logrus.Error(err)
		}
		// deploySpinner.UpdateText("Deploying node...\n")
		// Deploy
		_, deployErr := client.DeployNode(standFQDN, standUserData)
		if deployErr != nil {
			// deploySpinner.Fail("Failed to deploy node\n")
			logrus.Fatal(deployErr)
		}

		progressBar.Increment()
		progressBar.UpdateTitle("Waiting for deployment to complete and stand to be ready...")
		// deploySpinner.UpdateText("Waiting for deployment to complete and stand to be ready...\n")
		// Wait for deployment to complete
		deployWaitErr := client.WaitForReady(node)
		if deployWaitErr != nil {
			// deploySpinner.Fail("Failed to reach \"Ready\" state\n")
			logrus.Fatal(deployWaitErr)
		}

		pterm.Success.Println("Stand deployed successfully")
		progressBar.Increment()
		// deploySpinner.Success("Stand deployed successfully\n")
		// deploySpinner.Stop()
	},
}
