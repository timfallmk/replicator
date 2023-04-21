package provisioner

import (
	"encoding/base64"
	"errors"
	"os"
)

func encodeUserDataToBase64(userData []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(userData))
}

func readUserDataFromFile(filePath string) ([]byte, error) {
	file, openErr := os.ReadFile(filePath)
	if openErr != nil {
		return nil, openErr
	}
	return file, nil
}

func UserDataInputFromFile(filePath string) (string, error) {
	userData, readErr := readUserDataFromFile(filePath)
	if readErr != nil {
		return "", readErr
	}
	return encodeUserDataToBase64(userData), nil
}

func HostnameToFQDN(hostname string, pclient *ProvisionerClient) (string, error) {
	machineList, err := pclient.ListMachines()
	if err != nil {
		return "", err
	}

	for _, machine := range machineList {
		if machine.Hostname == hostname {
			return machine.FQDN, nil
		}
	}

	return "", errors.New("No matching node found with hostname: " + hostname)
}
