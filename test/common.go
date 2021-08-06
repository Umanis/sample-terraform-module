package test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type plan struct {
	Version string `json:"format_version"`
	Planned struct {
		Root struct {
			Child []childmodule `json:"child_modules"`
		} `json:"root_module"`
	} `json:"planned_values"`
}

type childmodule struct {
	Address   string     `json:"address"`
	Resources []resource `json:"resources"`
}

type resource struct {
	Address string                 `json:"address"`
	Values  map[string]interface{} `json:"values"`
}

func getAzureTerraformEnvVars() map[string]string {
	// getting enVars from environment variables
	ARM_CLIENT_ID := os.Getenv("AZURE_CLIENT_ID")
	ARM_CLIENT_SECRET := os.Getenv("AZURE_CLIENT_SECRET")
	ARM_SUBSCRIPTION_ID := os.Getenv("ARM_SUBSCRIPTION_ID")
	ARM_TENANT_ID := os.Getenv("AZURE_TENANT_ID")

	envVars := make(map[string]string)

	envVars["ARM_USE_MSI"] = "false"
	envVars["ARM_CLIENT_ID"] = ARM_CLIENT_ID
	envVars["ARM_CLIENT_SECRET"] = ARM_CLIENT_SECRET
	envVars["ARM_SUBSCRIPTION_ID"] = ARM_SUBSCRIPTION_ID
	envVars["ARM_TENANT_ID"] = ARM_TENANT_ID

	return envVars
}

func getPlannedResources(t *testing.T, jsonPlan string) map[string]resource {

	res := plan{}
	json.Unmarshal([]byte(jsonPlan), &res)

	actualLength := len(res.Planned.Root.Child)
	if actualLength != 1 {
		t.Fatal("No module")
	}

	actualModule := res.Planned.Root.Child[0]

	if len(actualModule.Resources) != 2 {
		t.Fatal("Wrong resource config")
	}

	buffer := make(map[string]resource)

	for i := 0; i < len(actualModule.Resources); i++ {
		resource := actualModule.Resources[i]
		buffer[resource.Address] = resource
	}

	return buffer
}

func generateVariables() map[string]interface{} {

	buffer := make(map[string]interface{})
	buffer["description"] = inputDescription
	buffer["location"] = inputLocation
	buffer["tags"] = make(map[string]string)
	buffer["custom_name"] = inputCustomName
	tags := make(map[string]string)
	tags[inputTagName] = inputTagValue

	buffer["tags"] = tags

	return buffer
}

var (
	inputCustomName  = "ZZZ00011la010"
	inputDescription = "descrip"
	inputLocation    = "France Central"
	inputTagName     = "tag"
	inputTagValue    = "tagvalue"
	outputLocation   = "francecentral"
)

func getResourceGroupsClient() resources.GroupsClient {
	resourcesClient := resources.NewGroupsClient(os.Getenv("ARM_SUBSCRIPTION_ID"))

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		log.Default().Print("OK")
		resourcesClient.Client.Authorizer = authorizer
	}
	if err != nil {
		log.Default().Print("KO")
		log.Default().Print(err)
	}

	return resourcesClient
}

func checkGroupExists(groupName string) bool {
	client := getResourceGroupsClient()

	_, err := client.CheckExistence(context.Background(), groupName)
	log.Default().Print(err)
	return err == nil
}
