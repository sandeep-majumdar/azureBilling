package cloudauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adeturner/azureBilling/observability"
)

type AzureServicePrincipalType struct {
	ClientId                       string `json:"clientId"`
	ClientSecret                   string `json:"ClientSecret"`
	SubscriptionId                 string `json:"SubscriptionId"`
	TenantId                       string `json:"TenantId"`
	ActiveDirectoryEndpointUrl     string `json:"ActiveDirectoryEndpointUrl"`
	ResourceManagerEndpointUrl     string `json:"ResourceManagerEndpointUrl"`
	ActiveDirectoryGraphResourceId string `json:"ActiveDirectoryGraphResourceId"`
	SqlManagementEndpointUrl       string `json:"SqlManagementEndpointUrl"`
	GalleryEndpointUrl             string `json:"GalleryEndpointUrl"`
	ManagementEndpointUrl          string `json:"ManagementEndpointUrl"`
}

func NewAzureServicePrincipalType() *AzureServicePrincipalType {
	return &AzureServicePrincipalType{}
}

func (spn *AzureServicePrincipalType) LoadFromFile() (err error) {

	var f *os.File
	var byteValue []byte

	file := os.Getenv("AZURE_AUTH_LOCATION")
	if file == "" {
		err = errors.New("Unable to find environment variable AZURE_AUTH_LOCATION")
	}

	if err == nil {
		f, err = os.Open(file)
		defer f.Close()

		if err == nil {

		} else {
			err = errors.New("Unable to access file described by env var AZURE_AUTH_LOCATION=%s" + file)
		}
	}

	if err == nil {
		byteValue, err = ioutil.ReadAll(f)
	}

	if err == nil {
		err = json.Unmarshal(byteValue, &spn)
	}

	if err != nil {
		observability.Error(err.Error())
	}

	return err
}

func (spn *AzureServicePrincipalType) GetClientId() string {
	return spn.ClientId
}

func (spn *AzureServicePrincipalType) GetClientSecret() string {
	return spn.ClientSecret
}

func (spn *AzureServicePrincipalType) GetTenant() string {
	return spn.TenantId

}

func (spn *AzureServicePrincipalType) Print() {
	observability.Info(fmt.Sprintf("%v", spn))
}
