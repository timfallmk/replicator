package main

import (
	"github.com/sirupsen/logrus"

	// maasclient "github.com/maas/gomaasclient/client"

	"ad.astra.com/gitlab/tss/replicator/cmd"
)

func main() {
	// Debug
	// logrus.SetLevel(logrus.DebugLevel)
	// Info
	logrus.SetLevel(logrus.InfoLevel)
	err := cmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
	// client, err := maasclient.GetClient("http://empok-nor:5240/MAAS", "VDwajrSanHxEnFZUKQ:b85sQJqc8qdm5CgHRP:Z8ycUq9ZaunWNGAGcHtEFe3SycqNuGHh", "2.0")
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// // Get all nodes
	// nodes, err := client.Machines.Get()
	// if err != nil {
	// 	logrus.Fatal(err)
	// }

	// // Print all nodes
	// for _, node := range nodes {
	// 	if node.Pool.Name == "Stands" {
	// 		logrus.Infof("+%v", node.IPAddresses)
	// 	}
	// }
}
