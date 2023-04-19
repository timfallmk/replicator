package provisioner

import (
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	maasEntity "github.com/maas/gomaasclient/entity"
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

func (pclient *ProvisionerClient) findNewMachines() ([]maasEntity.Machine, error) {
	// Get all nodes
	nodes, err := pclient.Client.Machines.Get()
	switch {
	case err != nil:
		return []maasEntity.Machine{}, err
	case len(nodes) == 0:
		return []maasEntity.Machine{}, fmt.Errorf("No nodes found")
	}

	// logrus.Debugf("Nodes: %+v", nodes)

	// Find new nodes
	newMachines := []maasEntity.Machine{}
	for _, node := range nodes {
		switch {
		// Find nodes matching the "New" status
		case node.StatusName == "New":
			newMachines = append(newMachines, node)
			// TODO: There needs to be more matching conditions. Age?
		}
	}

	return newMachines, nil
}

func (pclient *ProvisionerClient) CommissionNode(hostname string) (*maasEntity.Machine, error) {
	// Get systemID from hostname
	systemID, err := pclient.nodeHostnameToSystemID(hostname)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("SystemID: %v", systemID)

	comissionParameters := maasEntity.MachineCommissionParams{}

	// Commission node
	machine, err := pclient.Client.Machine.Commission(systemID, &comissionParameters)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

func (pclient *ProvisionerClient) isComissioned(machine *maasEntity.Machine) bool {
	return machine.CommissioningStatus == 1
}

func (pclient *ProvisionerClient) isReady(machine *maasEntity.Machine) bool {
	return machine.StatusName == "Ready"
}

func (pclient *ProvisionerClient) WaitForComissioned(machine *maasEntity.Machine) error {
	// Wait for node to be comissioned
	comissioned := make(chan bool, 1)
	comissioned <- false
	go func() {
		if pclient.isComissioned(machine) {
			comissioned <- true
			return
		}
		time.Sleep(1 * time.Second)
	}()

	select {
	case <-comissioned:
		return nil
	case <-time.After(10 * time.Minute):
		return fmt.Errorf("Timed out waiting for node to be comissioned")
	}
}

func (pclient *ProvisionerClient) WaitForReady(machine *maasEntity.Machine) error {
	// Wait for node to be ready
	ready := make(chan bool, 1)
	ready <- false
	go func() {
		if pclient.isReady(machine) {
			ready <- true
			return
		}
		time.Sleep(1 * time.Second)
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(10 * time.Minute):
		return fmt.Errorf("Timed out waiting for node to be ready")
	}
}

func (pclient *ProvisionerClient) allocateToStands(machine *maasEntity.Machine) error {
	// Allocate machine to pool "Stands"
	_, err := pclient.Client.Machine.Update(machine.SystemID, &maasEntity.MachineParams{
		Pool: "Stands",
	}, map[string]string{})
	return err
}

func (pclient *ProvisionerClient) DeployNode(hostname string, userData string) (*maasEntity.Machine, error) {
	// Get systemID from hostname
	systemID, err := pclient.nodeHostnameToSystemID(hostname)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("SystemID: %v", systemID)

	deployParameters := maasEntity.MachineDeployParams{
		// Base64 encoded user-data
		UserData:     userData,
		DistroSeries: "jammy",
		HWEKernel:    "hwe-22.04-lowlatency-edge",
		// TODO: Hadware Sync is missing from the API library
	}

	// Deploy node
	machine, err := pclient.Client.Machine.Deploy(systemID, &deployParameters)
	if err != nil {
		return nil, err
	}

	return machine, nil
}
