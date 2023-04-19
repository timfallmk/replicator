package provisioner

import (
	"encoding/base64"
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
