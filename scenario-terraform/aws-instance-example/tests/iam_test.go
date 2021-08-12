package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	testStructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/likexian/gokit/assert"
	"github.com/tidwall/gjson"
)

/**
1. Stick to the naming module rules, (group, user, etc)
2. Check password_reset_required == true
3. Check attach_iam_self_management_policy = false
4. Check if user access to all resource EC2
**/

/** Test IAM Role Create User
	- Check naming convention of the module "<name>-gdp"
	- Password change == true
	- Iam management self policy == false
**/
func TestUT_IAMRoleCreateUser(t *testing.T) {
	t.Parallel()

	// IAM terraform folder
	iamFolder := testStructure.CopyTerraformFolderToTemp(t, "../../", "aws-instance-example/iam")

	awsRegion := "ap-southeast-1"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: iamFolder,

		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	// Destroy
	defer terraform.Destroy(t, terraformOptions)

	// Init and Plan
	terraform.InitAndPlan(t, terraformOptions)

	tfPlanOutput := "terraform.tfplan"
	terraform.RunTerraformCommand(t, terraformOptions, "plan", "-out="+tfPlanOutput)

	planJson, _ := terraform.RunTerraformCommandAndGetStdoutE(t, terraformOptions, "show", "-json", tfPlanOutput)

	strJson, _ := json.Marshal(planJson)
	fmt.Print(string(strJson))
	userName := gjson.Get(string(strJson), "planned_values.root_module.child_modules[0].resources[0].values.name")

	fmt.Printf(userName.String())
	// Check if username created contains gdp substring

	assert.Contains(t, userName, "-gdp")

	// Check if password_reset_required == true

	// Check if attach_iam_self_management_policy = false

}

// Test IAM Role EC2 Instance
func TestUT_IAMRoleEC2(t *testing.T) {
	t.Parallel()
	// Check if user access to all resource EC2

}
