package cmd

import (
	"net/url"

	// TODO: Remove this dependency. Hack for terminal wrapping.
	// prettyText "github.com/jedib0t/go-pretty/v6/text"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"ad.astra.com/gitlab/tss/replicator/pkg/provisioner"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stand",
	Long:  `Delete a stand. This will delete a stand and remove it from the MAAS server.`,
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

		pterm.DefaultBigText.WithLetters(putils.LettersFromStringWithStyle("WARNING", pterm.FgRed.ToStyle())).Render()
		pterm.DefaultBasicText.Println("This will delete the stand and remove it completely from the MAAS server!")
		pterm.DefaultBasicText.Println("If you want to only erase the machine and return it to it's default state\n use the 'erase' command instead.")

		confimed, err := pterm.DefaultInteractiveContinue.WithOptions(
			[]string{"continue", "abort"},
		).WithDefaultValue(
			"abort",
		).Show()
		if err != nil {
			logrus.Fatal(err)
		}

		switch confimed {
		case "Continue":
			break
		case "Abort":
			pterm.Println("Aborted")
			logrus.Exit(0)
		}

		spinnerDelete, err := pterm.DefaultSpinner.Start("Deleting machine...")
		if err != nil {
			logrus.Error(err)
		}

		delErr := pclient.RemoveMachine(args[0])
		if delErr != nil {
			spinnerDelete.Fail("Failed to delete machine after ", spinnerDelete.Delay)
			logrus.Fatal(delErr)
		}

		spinnerDelete.Success("Deleted machine after ", spinnerDelete.Delay)
	},
}
