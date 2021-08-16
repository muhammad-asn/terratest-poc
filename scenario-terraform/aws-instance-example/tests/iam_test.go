package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	testStructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

/**
1. Stick to the naming module rules, (group, user, etc)
2. Check password_reset_required == true
3. Check attach_iam_self_management_policy = false
https://github.com/gruntwork-io/terratest/blob/master/test/terraform_aws_example_plan_test.go
**/
type Tags struct {
	Organization    string
	Project         string
	Terraform       string
	Email           string
	BastionUserName string
}

var awsRegion = "ap-southeast-1"

var uTags = Tags{
	Organization:    "bonek",
	Project:         "bonek",
	Terraform:       "true",
	Email:           "@mania.com",
	BastionUserName: "-mania"}

func GetTerraformOptions(t *testing.T) *terraform.Options {

	// IAM terraform folder and plan
	var iamFolder = testStructure.CopyTerraformFolderToTemp(t, "../../", "aws-instance-example/iam")
	var planFilePath = filepath.Join(iamFolder, "plan.out")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: iamFolder,

		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
		PlanFilePath: planFilePath,
	})

	return terraformOptions

}

func UT_IAMRoleCreateUser(t *testing.T, terraformOptions *terraform.Options, plan *terraform.PlanStruct) {

	resChildModules := plan.RawPlan.PlannedValues.RootModule.ChildModules

	for i := 0; i < len(resChildModules); i++ {

		userName := resChildModules[i].Address

		// Resource user IAM -> even numbers
		if strings.Contains(userName, "mania") {

			awsIamUserPlan := resChildModules[i].Address + ".aws_iam_user.this[0]"
			awsIamProfilePlan := resChildModules[i].Address + ".aws_iam_user_login_profile.this[0]"

			terraform.RequirePlannedValuesMapKeyExists(t, plan, awsIamUserPlan)
			userResource := plan.ResourcePlannedValuesMap[awsIamUserPlan]
			userProfile := plan.ResourcePlannedValuesMap[awsIamProfilePlan]
			userTags := userResource.AttributeValues["tags"].(map[string]interface{})

			// Check if contains specified tags
			assert.Contains(t, userTags["Project"], uTags.Project, userTags["Project"])
			assert.Contains(t, userTags["Organization"], uTags.Organization, userTags["Organization"])
			assert.Contains(t, userTags["Terraform"], uTags.Terraform, userTags["Terraform"])
			assert.Contains(t, userTags["Email"], uTags.Email, userTags["Email"])
			assert.Contains(t, userTags["BastionUserName"], uTags.BastionUserName, userTags["BastionUserName"])

			// Check if username created contains mania substring
			userName := userResource.AttributeValues["name"]
			assert.Contains(t, userName, "-mania")

			// Check destroy user true
			destroyUser := userResource.AttributeValues["force_destroy"]

			assert.Equal(t, destroyUser, true)

			// Check if password_reset_required == true
			pwdResetReq := userProfile.AttributeValues["password_reset_required"]
			assert.Equal(t, pwdResetReq, true)

			// Check if create_iam_access_key == false
			createIamAccessKey := plan.RawPlan.Config.RootModule.ModuleCalls[userName.(string)].Expressions["create_iam_access_key"].ConstantValue
			assert.Equal(t, createIamAccessKey, false)
		} else {
			continue
		}

	}

}

func UT_IAMCreatePolicy(t *testing.T, terraformOptions *terraform.Options, plan *terraform.PlanStruct) {

	resChildModules := plan.RawPlan.PlannedValues.RootModule.ChildModules

	for i := 0; i < len(resChildModules); i++ {

		moduleName := resChildModules[i].Address

		// module.<resource_name> -> <resource_name>
		expectedResourceName := strings.ReplaceAll(moduleName[7:], "_", "-")

		if strings.Contains(moduleName, "module.policy") {

			resourceAddress := moduleName + ".aws_iam_policy.policy"
			terraform.RequirePlannedValuesMapKeyExists(t, plan, resourceAddress)
			policyResource := plan.ResourcePlannedValuesMap[resourceAddress]

			gotResourceName := policyResource.AttributeValues["name"].(string)

			assert.Equal(t, gotResourceName, expectedResourceName)
		}
	}
}

func UT_IAMCreateGroupPolicy(t *testing.T, terraformOptions *terraform.Options, plan *terraform.PlanStruct) {

	resChildModules := plan.RawPlan.PlannedValues.RootModule.ChildModules

	for i := 0; i < len(resChildModules); i++ {

		moduleName := resChildModules[i].Address

		// module.<resource_name> -> <resource_name>
		expectedResourceName := strings.ReplaceAll(moduleName[7:], "_", "-")

		if strings.Contains(moduleName, "module.group_policy") {

			resourceAddress := moduleName + ".aws_iam_group_membership.this[0]"
			terraform.RequirePlannedValuesMapKeyExists(t, plan, resourceAddress)
			policyResource := plan.ResourcePlannedValuesMap[resourceAddress]

			gotResourceName := policyResource.AttributeValues["name"].(string)

			assert.Equal(t, gotResourceName, expectedResourceName)
		}
	}

}

func TestUT_Start(t *testing.T) {

	t.Parallel()

	terraformOptions := GetTerraformOptions(t)

	// Destroy
	defer terraform.Destroy(t, terraformOptions)

	plan := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

	UT_IAMRoleCreateUser(t, terraformOptions, plan)
	UT_IAMCreatePolicy(t, terraformOptions, plan)
	UT_IAMCreateGroupPolicy(t, terraformOptions, plan)

}
