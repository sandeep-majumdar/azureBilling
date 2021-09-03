package cloudauth

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type AzureSession struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer
}

func readJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err // errors.Wrap(err, "Can't open the file")
	}

	contents := make(map[string]interface{})
	err = json.Unmarshal(data, &contents)

	//if err != nil {
	//err = err // errors.Wrap(err, "Can't unmarshal file")
	//}

	return &contents, err
}

func NewSessionFromFile() (*AzureSession, error) {
	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)

	if err != nil {
		return nil, err
	}

	authInfo, err := readJSON(os.Getenv("AZURE_AUTH_LOCATION"))

	if err != nil {
		return nil, err
	}

	sess := &AzureSession{
		SubscriptionID: (*authInfo)["subscriptionId"].(string),
		Authorizer:     authorizer,
	}

	return sess, nil
}
