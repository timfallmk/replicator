package cmd

import (
	"bytes"
	"net/url"
	"sort"

	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"ad.astra.com/gitlab/tss/replicator/pkg/provisioner"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all the stands",
	Long:    `List all the currently known stands and their status.`,
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		pURL, err := url.Parse(provisionerURL)
		if err != nil {
			logrus.Fatal(err)
		}
		pclient, err := provisioner.New(*pURL, provisionerToken)
		if err != nil {
			logrus.Fatal(err)
		}

		spinnerList, err := pterm.DefaultSpinner.Start("Getting list of machines...")
		if err != nil {
			logrus.Error(err)
		}

		machineList, err := pclient.ListMachines()
		if err != nil {
			spinnerList.Fail("Failed to get list of machines.")
			logrus.Fatal(err)
		}

		// TODO: This is a hack to get the table to display properly. Clean this up.
		table := pterm.DefaultTable.WithBoxed().WithHasHeader().WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(pterm.TableData{
			{"Hostname 📡", "FQDN 📇", "IP 🛜", "Pool 🏊‍♀️", "Power 🔌", "Status 👍"},
		})
		for _, machine := range machineList {
			// Make sure the IP addresses are in order
			sort.Slice(machine.IP, func(i, j int) bool {
				return bytes.Compare(machine.IP[i], machine.IP[j]) < 0
			})

			var ipRange string
			for _, ip := range machine.IP {
				ipRange += ip.String() + "\n\n"
			}
			table.Data = append(table.Data, []string{machine.Hostname, machine.FQDN, ipRange, machine.Pool, machine.Power, machine.Status})
		}
		spinnerList.Success("Got list of machines after ", spinnerList.Delay)
		rendErr := table.Render()
		if err != nil {
			logrus.Error(rendErr)
		}
	},
}
