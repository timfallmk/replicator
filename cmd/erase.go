package cmd

import (
	"net/url"

	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"ad.astra.com/gitlab/tss/replicator/pkg/provisioner"
)

var eraseCmd = &cobra.Command{
	Use:   "erase",
	Short: "Erase a stand",
	Long:  `Erase a stand. This will erase a stand and return it to it's default state.`,
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

		pterm.Warning.Printf("Erasing machine %s\n", args[0])

		continued, err := pterm.DefaultInteractiveContinue.WithOptions(
			[]string{"continue", "abort"},
		).WithDefaultValue(
			"abort",
		).Show()
		if err != nil {
			logrus.Fatal(err)
		}

		switch continued {
		case "continue":
			break
		case "abort":
			pterm.Println("Aborted")
			logrus.Exit(0)
		}

		spinnerErase, err := pterm.DefaultSpinner.Start("Erasing machine...")
		if err != nil {
			logrus.Error(err)
		}

		eraseErr := pclient.EraseMachine(args[0])
		if eraseErr != nil {
			spinnerErase.Fail("Failed to erase machine after ", spinnerErase.Delay)
			logrus.Fatal(eraseErr)
		}

		spinnerErase.Success("Erased machine after ", spinnerErase.Delay)
	},
}
