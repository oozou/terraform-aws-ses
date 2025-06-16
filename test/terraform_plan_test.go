package test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformPlan(t *testing.T) {
	t.Parallel()

	t.Run("EmailVerificationPlan", func(t *testing.T) {
		testEmailVerificationPlan(t)
	})

	t.Run("DomainVerificationPlan", func(t *testing.T) {
		testDomainVerificationPlan(t)
	})
}

func testEmailVerificationPlan(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/email-verification",
		Vars: map[string]interface{}{
			"prefix":       "test",
			"environment":  "dev",
			"tags":         map[string]string{"test": "true"},
		},
		PlanFilePath: "./email-plan.tfplan",
	}

	// Generate terraform plan
	terraform.Init(t, terraformOptions)
	planOutput := terraform.Plan(t, terraformOptions)

	// Validate that SES email identity will be created
	assert.Contains(t, planOutput, "aws_ses_email_identity.this", "Plan should include SES email identity creation")
	assert.Contains(t, planOutput, "will be created", "Resources should be marked for creation")

	// Parse plan as JSON for detailed validation
	planJson := terraform.Show(t, terraformOptions)
	
	// Validate SES email identity configuration
	validateEmailIdentityInPlan(t, planJson)
}

func testDomainVerificationPlan(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/domain-verification",
		Vars: map[string]interface{}{
			"prefix":       "test",
			"environment":  "dev",
			"tags":         map[string]string{"test": "true"},
		},
		PlanFilePath: "./domain-plan.tfplan",
	}

	// Generate terraform plan
	terraform.Init(t, terraformOptions)
	planOutput := terraform.Plan(t, terraformOptions)

	// Validate that SES domain identity will be created
	assert.Contains(t, planOutput, "aws_ses_domain_identity.this", "Plan should include SES domain identity creation")
	
	// Validate DMARC record will be created
	assert.Contains(t, planOutput, "aws_route53_record.dmarc", "Plan should include DMARC record creation")
	
	// Validate DKIM configuration
	assert.Contains(t, planOutput, "aws_ses_domain_dkim.this", "Plan should include SES domain DKIM")

	// Parse plan as JSON for detailed validation
	planJson := terraform.Show(t, terraformOptions)
	
	// Validate domain identity, DKIM, and DMARC configuration
	validateDomainIdentityInPlan(t, planJson)
	validateDMARCInPlan(t, planJson)
}

func validateEmailIdentityInPlan(t *testing.T, planJson string) {
	var plan map[string]interface{}
	err := json.Unmarshal([]byte(planJson), &plan)
	require.NoError(t, err, "Failed to parse plan JSON")

	// Navigate to planned values
	plannedValues, ok := plan["planned_values"].(map[string]interface{})
	require.True(t, ok, "Plan should have planned_values")

	rootModule, ok := plannedValues["root_module"].(map[string]interface{})
	require.True(t, ok, "Plan should have root_module")

	childModules, ok := rootModule["child_modules"].([]interface{})
	require.True(t, ok && len(childModules) > 0, "Plan should have child_modules")

	// Find the SES module
	sesModule := findModuleByAddress(childModules, "module.ses")
	require.NotNil(t, sesModule, "SES module should exist in plan")

	resources, ok := sesModule["resources"].([]interface{})
	require.True(t, ok, "SES module should have resources")

	// Look for SES email identity resource
	emailIdentityFound := false
	for _, resource := range resources {
		if res, ok := resource.(map[string]interface{}); ok {
			if resType, ok := res["type"].(string); ok && resType == "aws_ses_email_identity" {
				emailIdentityFound = true
				
				// Validate email identity configuration
				values, ok := res["values"].(map[string]interface{})
				require.True(t, ok, "Email identity should have values")
				
				email, ok := values["email"].(string)
				require.True(t, ok, "Email identity should have email value")
				assert.Contains(t, email, "@", "Email should be a valid email format")
				
				break
			}
		}
	}
	
	assert.True(t, emailIdentityFound, "SES email identity should be found in plan")
}

func validateDomainIdentityInPlan(t *testing.T, planJson string) {
	var plan map[string]interface{}
	err := json.Unmarshal([]byte(planJson), &plan)
	require.NoError(t, err, "Failed to parse plan JSON")

	// Navigate to planned values
	plannedValues, ok := plan["planned_values"].(map[string]interface{})
	require.True(t, ok, "Plan should have planned_values")

	rootModule, ok := plannedValues["root_module"].(map[string]interface{})
	require.True(t, ok, "Plan should have root_module")

	childModules, ok := rootModule["child_modules"].([]interface{})
	require.True(t, ok && len(childModules) > 0, "Plan should have child_modules")

	// Find the SES module
	sesModule := findModuleByAddress(childModules, "module.ses")
	require.NotNil(t, sesModule, "SES module should exist in plan")

	resources, ok := sesModule["resources"].([]interface{})
	require.True(t, ok, "SES module should have resources")

	// Look for SES domain identity and DKIM resources
	domainIdentityFound := false
	dkimFound := false
	
	for _, resource := range resources {
		if res, ok := resource.(map[string]interface{}); ok {
			if resType, ok := res["type"].(string); ok {
				switch resType {
				case "aws_ses_domain_identity":
					domainIdentityFound = true
					
					// Validate domain identity configuration
					values, ok := res["values"].(map[string]interface{})
					require.True(t, ok, "Domain identity should have values")
					
					domain, ok := values["domain"].(string)
					require.True(t, ok, "Domain identity should have domain value")
					assert.True(t, len(domain) > 0, "Domain should not be empty")

				case "aws_ses_domain_dkim":
					dkimFound = true
					
					// Validate DKIM configuration
					values, ok := res["values"].(map[string]interface{})
					require.True(t, ok, "DKIM should have values")
					
					domain, ok := values["domain"].(string)
					require.True(t, ok, "DKIM should have domain value")
					assert.True(t, len(domain) > 0, "DKIM domain should not be empty")
				}
			}
		}
	}
	
	assert.True(t, domainIdentityFound, "SES domain identity should be found in plan")
	assert.True(t, dkimFound, "SES domain DKIM should be found in plan")
}

func validateDMARCInPlan(t *testing.T, planJson string) {
	var plan map[string]interface{}
	err := json.Unmarshal([]byte(planJson), &plan)
	require.NoError(t, err, "Failed to parse plan JSON")

	// Navigate to planned values
	plannedValues, ok := plan["planned_values"].(map[string]interface{})
	require.True(t, ok, "Plan should have planned_values")

	rootModule, ok := plannedValues["root_module"].(map[string]interface{})
	require.True(t, ok, "Plan should have root_module")

	childModules, ok := rootModule["child_modules"].([]interface{})
	require.True(t, ok && len(childModules) > 0, "Plan should have child_modules")

	// Find the SES module
	sesModule := findModuleByAddress(childModules, "module.ses")
	require.NotNil(t, sesModule, "SES module should exist in plan")

	resources, ok := sesModule["resources"].([]interface{})
	require.True(t, ok, "SES module should have resources")

	// Look for DMARC Route53 record
	dmarcRecordFound := false
	
	for _, resource := range resources {
		if res, ok := resource.(map[string]interface{}); ok {
			if resType, ok := res["type"].(string); ok && resType == "aws_route53_record" {
				// Check if this is the DMARC record
				if name, ok := res["name"].(string); ok && strings.Contains(name, "dmarc") {
					dmarcRecordFound = true
					
					// Validate DMARC record configuration
					values, ok := res["values"].(map[string]interface{})
					require.True(t, ok, "DMARC record should have values")
					
					recordType, ok := values["type"].(string)
					require.True(t, ok, "DMARC record should have type")
					assert.Equal(t, "TXT", recordType, "DMARC record should be TXT type")
					
					// Check if record name contains _dmarc
					recordName, ok := values["name"].(string)
					require.True(t, ok, "DMARC record should have name")
					assert.Contains(t, recordName, "_dmarc", "DMARC record name should contain _dmarc")
					
					// Check records array
					if records, ok := values["records"].([]interface{}); ok && len(records) > 0 {
						if firstRecord, ok := records[0].(string); ok {
							assert.Contains(t, firstRecord, "v=DMARC1", "DMARC record should contain v=DMARC1")
						}
					}
					
					break
				}
			}
		}
	}
	
	assert.True(t, dmarcRecordFound, "DMARC Route53 record should be found in plan")
}

func findModuleByAddress(modules []interface{}, address string) map[string]interface{} {
	for _, module := range modules {
		if mod, ok := module.(map[string]interface{}); ok {
			if addr, ok := mod["address"].(string); ok && addr == address {
				return mod
			}
		}
	}
	return nil
}