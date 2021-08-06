package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestPlanCommonValues(t *testing.T) {
	t.Parallel()

	terraformVars := generateVariables()

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./common_values_fixture",
		Vars:         terraformVars,
		EnvVars:      getAzureTerraformEnvVars(),
		NoColor:      true,
	})

	tfPlanOutput := "terraform.tfplan"

	terraform.Init(t, tfOptions)
	terraform.RunTerraformCommand(t, tfOptions, terraform.FormatArgs(tfOptions, "plan", "-out="+tfPlanOutput)...)

	tfOptionsEmpty := &terraform.Options{}
	planJSON, err := terraform.RunTerraformCommandAndGetStdoutE(
		t, tfOptions, terraform.FormatArgs(tfOptionsEmpty, "show", "-json", tfPlanOutput)...,
	)

	if err != nil {
		t.Fatal(err)
	}

	planned := getPlannedResources(t, planJSON)

	groupResource := planned["module.group.azurerm_resource_group.self"]

	// Location
	actualLocation := groupResource.Values["location"].(string)
	if actualLocation != outputLocation {
		t.Fatal("Wrong location")
	}

	// tags
	actualTags := groupResource.Values["tags"].(map[string]interface{})
	def, exists := actualTags[inputTagName]
	if !exists || def != inputTagValue {
		t.Fatal("Wrong tag definition")
	}
	descrip, exists := actualTags["description"]
	if !exists || descrip != inputDescription {
		t.Fatal("Wrong description")
	}
}
