package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestCustomNameValues(t *testing.T) {
	t.Parallel()

	terraformVars := generateVariables()

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./custom_name_fixture/",
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

	// name
	actualName := groupResource.Values["name"].(string)
	if actualName != inputCustomName {
		t.Fatal("Wrong name")
	}
}
