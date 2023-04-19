package provisioner

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type MachineList struct {
	Hostname string
	FQDN     string
	IP       []net.IP
	Pool     string
	Power    string
	Status   string
}

type MachineDetails struct {
	Hostname          string
	FQDN              string
	Domain            string
	IP                []net.IP
	Pool              string
	Power             string
	Status            string
	NetworkInterfaces []string
	NetBoot           bool
	Arch              string
	OS                string
	Kernel            string
	PowerType         string
	StatusMessage     string
}

func (pclient *ProvisionerClient) ListMachines() ([]MachineList, error) {
	// Get all nodes
	nodes, err := pclient.Client.Machines.Get()
	if err != nil {
		return []MachineList{}, err
	}

	// Rerturn list of nodes
	list := []MachineList{}
	for _, node := range nodes {
		list = append(list, MachineList{
			Hostname: node.Hostname,
			FQDN:     node.FQDN,
			IP:       node.IPAddresses,
			Pool:     node.Pool.Name,
			Power:    node.PowerState,
			Status:   node.StatusName,
		})
	}
	return list, nil

	// // Print all nodes
	// for _, node := range nodes {
	// 	if node.Pool.Name == "Stands" {
	// 		logrus.Printf("+%v", node.IPAddresses)
	// 	}
	// }
}

func (pclient *ProvisionerClient) GetMachineDetails(hostname string) (MachineDetails, error) {
	// Get systemID from hostname
	systemID, err := pclient.nodeHostnameToSystemID(hostname)
	if err != nil {
		return MachineDetails{}, err
	}
	logrus.Debugf("SystemID: %v", systemID)

	// Get node details
	node, err := pclient.Client.Machine.Get(systemID)
	if err != nil {
		return MachineDetails{}, err
	}
	logrus.Debugf("Node: %+v", node)

	// Rerturn list of nodes

	details := MachineDetails{
		Hostname:          node.Hostname,
		FQDN:              node.FQDN,
		Domain:            node.Domain.Name,
		IP:                node.IPAddresses,
		Pool:              node.Pool.Name,
		Power:             node.PowerState,
		Status:            node.StatusName,
		NetworkInterfaces: node.BootInterface.Children,
		NetBoot:           node.Netboot,
		Arch:              node.Architecture,
		OS:                node.OSystem,
		PowerType:         node.PowerType,
		StatusMessage:     node.StatusMessage,
	}
	return details, nil
}

func (pclient *ProvisionerClient) nodeHostnameToSystemID(hostname string) (string, error) {
	// Get all nodes
	nodes, err := pclient.Client.Machines.Get()
	switch {
	case err != nil:
		return "", err
	case len(nodes) == 0:
		return "", fmt.Errorf("No nodes found")
	}

	// logrus.Debugf("Nodes: %+v", nodes)

	// Find node by hostname
	for _, node := range nodes {
		switch {
		case node.Hostname == hostname:
			logrus.Debugf("Node Name: %s", node.Hostname)
			logrus.Debugf("Node ID: %s", node.SystemID)
			return node.SystemID, nil
		case node.FQDN == hostname:
			logrus.Debugf("Node Name: %s", node.FQDN)
			logrus.Debugf("Node ID: %s", node.SystemID)
			return node.SystemID, nil
		}
	}

	return "", fmt.Errorf("No matching node found with hostname: %s", hostname)
}
