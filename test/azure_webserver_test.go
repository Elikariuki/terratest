package test

import (
	"time"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "625993dd-b03e-408c-bdd0-dabcd4100720"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "kari0117",
		},
	}

	defer func() {
        // Destroy logic
        err := retryDestroy(t, terraformOptions, 15*time.Minute, 30*time.Second) // Retry for 15 minutes
        if err != nil {
            t.Errorf("Error during destroy: %s", err)
        }
    }()
 
    // Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
    terraform.InitAndApply(t, terraformOptions)
 
    // Run `terraform output` to get the value of output variable
    vmName := terraform.Output(t, terraformOptions, "vm_name")
    resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
    nicName := terraform.Output(t, terraformOptions, "nic_name")
    nicID := terraform.Output(t, terraformOptions, "nic_id")
    vmNicIds := terraform.Output(t, terraformOptions, "vm_nic")
    vmVersion := terraform.Output(t, terraformOptions, "vm_version")
    vmSKU := terraform.Output(t, terraformOptions, "vm_sku")
 
    // Confirm VM exists
    assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
 
    // - Validate that the NIC exists
    assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))
 
    // - Validate that the NIC is connected to the VM
    assert.Contains(t, vmNicIds, nicID)
 
    // - Validate vm version and sku
    assert.Equal(t, "latest", vmVersion)
    assert.Equal(t, "22_04-lts-gen2", vmSKU)
 
}
 
// retry function for destroy with error filtering
func retryDestroy(t *testing.T, terraformOptions *terraform.Options, maxDuration time.Duration, waitBetween time.Duration) error {
    startTime := time.Now()
    var err error
    for {
        // Capture both returned values: result and error
        _, err = terraform.DestroyE(t, terraformOptions)
        if err == nil {
            // Success: No error, we are done
            return nil
        }
 
        // Check if we have exceeded the maxDuration
        if time.Since(startTime) > maxDuration {
            return err
        }
 
        // Check if the error is related to "NetworkSecurityGroupOldReferencesNotCleanedUp"
        if strings.Contains(err.Error(), "NetworkSecurityGroupOldReferencesNotCleanedUp") {
            t.Logf("Network Security Group still has references. Retrying destroy: %v", err)
            time.Sleep(waitBetween) // Wait before retrying
        } else if strings.Contains(err.Error(), "InternalServerError") {
            // Handle InternalServerError as before
            t.Logf("Error during destroy: %v. Retrying...", err)
            time.Sleep(waitBetween) // Wait before retrying
        } else {
            // If it's a different error, return it immediately
            return err
        }
    }
}
