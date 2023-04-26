package cmd

import (
	"bytes"
	"net/url"
	"sort"
	"strconv"

	// TODO: Remove this dependency. Hack for terminal wrapping.
	prettyText "github.com/jedib0t/go-pretty/v6/text"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"ad.astra.com/gitlab/tss/replicator/pkg/provisioner"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get stand details",
	Long:  `Get details for a specific stand.`,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		pURL, err := url.Parse(provisionerURL)
		if err != nil {
			logrus.Fatal(err)
		}
		pclient, err := provisioner.New(*pURL, provisionerToken)
		if err != nil {
			logrus.Fatal(err)
		}

		spinnerGet, err := pterm.DefaultSpinner.Start("Getting machine details...")
		if err != nil {
			logrus.Error(err)
		}

		details, detailErr := pclient.GetMachineDetails(args[0])
		if detailErr != nil {
			spinnerGet.Fail("Failed to get machine details after ", spinnerGet.Delay)
			logrus.Fatal(detailErr)
		}

		// TODO: This is a hack to get the table to display properly. Clean this up.
		table := pterm.DefaultTable.WithBoxed().WithHasHeader().WithHasHeader().WithRowSeparator("-").WithHeaderRowSeparator("-").WithData(pterm.TableData{
			{"Hostname", "ID", "FQDN", "IP", "Pool", "Power", "Status", "NetworkInterfaces", "NetBoot", "Arch", "OS", "Kernel", "Power Type", "Status Message"},
		})

		// Make sure the IP addresses are in order
		sort.Slice(details.IP, func(i, j int) bool {
			return bytes.Compare(details.IP[i], details.IP[j]) < 0
		})

		var ipRange string
		for _, ip := range details.IP {
			ipRange += ip.String() + "\n\n"
		}

		var netInterfaces string
		for _, netI := range details.NetworkInterfaces {
			netInterfaces += netI + "\n\n"
		}

		table.Data = append(table.Data, []string{details.Hostname, details.SystemID, details.FQDN, ipRange, details.Pool, details.Power, details.Status, netInterfaces, strconv.FormatBool(details.NetBoot), details.Arch, details.OS, details.Kernel, details.PowerType, prettyText.WrapSoft(details.StatusMessage, 30)})
		spinnerGet.Success("Got got details after ", spinnerGet.Delay)
		rendErr := table.Render()
		if err != nil {
			logrus.Error(rendErr)
		}
	},
}
