package provisioner

import (
	"net/url"

	//"github.com/sirupsen/logrus"

	maasclient "github.com/maas/gomaasclient/client"
)

type ProvisionerClient struct {
	Client *maasclient.Client
}

func New(provisionerURL url.URL, provisionerToken string) (*ProvisionerClient, error) {
	client, err := maasclient.GetClient(provisionerURL.String(), provisionerToken, "2.0")
	if err != nil {
		return nil, err
	}

	return &ProvisionerClient{
		Client: client,
	}, nil

}
