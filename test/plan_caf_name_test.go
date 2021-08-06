package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestCafNameValues(t *testing.T) {
	t.Parallel()

	terraformVars := generateVariables()

	terraformVars["custom_name"] = ""
	terraformVars["caf_prefixes"] = []string{"a010", "abc", "6666", "l", "02"}

	outputName := "a010-abc-6666-l-02"

	tfOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./caf_name_fixture/",
		Vars:         terraformVars,
		EnvVars:      getAzureTerraformEnvVars(),
		NoColor:      true,
	})

	defer terraform.Destroy(t, tfOptions)

	terraform.InitAndApply(t, tfOptions)

	assert.Equal(t, true, checkGroupExists(outputName))

}
