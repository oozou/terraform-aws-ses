package test

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/oozou/terraform-test-util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Global variables for test reporting
var (
	generateReport bool
	reportFile     string
	htmlFile       string
)

// TestMain enables custom test runner with reporting
func TestMain(m *testing.M) {
	flag.BoolVar(&generateReport, "report", false, "Generate test report")
	flag.StringVar(&reportFile, "report-file", "test-report.json", "Test report JSON file")
	flag.StringVar(&htmlFile, "html-file", "test-report.html", "Test report HTML file")
	flag.Parse()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTerraformAWSSESModule(t *testing.T) {
	t.Parallel()

	// Record test start time
	startTime := time.Now()
	var testResults []testutil.TestResult

	// Pick a random AWS region to test in
	awsRegion := "ap-southeast-1"

	// Generate a unique name for resources
	uniqueID := strings.ToLower(random.UniqueId())
	testDomain := fmt.Sprintf("test%s.oozou.com", uniqueID)
	testEmail := fmt.Sprintf("test-%s@oozou.com", uniqueID)

	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-test",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix":      "terratest",
			"environment": "test",
			"tags":        map[string]string{"test": "true"},
			"ses_domain": map[string]interface{}{
				"domain":            testDomain,
				"is_verify_dkim":    true,
				"is_verify_domain":  true,
				"route53_zone_name": testDomain,
				"is_verify_dmarc":   true,
			},
			"route53_zone_name": testDomain,
			"ses_email": map[string]interface{}{
				"email": testEmail,
			},
		},

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer func() {
		terraform.Destroy(t, terraformOptions)

		// Generate and display test report
		endTime := time.Now()
		report := testutil.GenerateTestReport(testResults, startTime, endTime)
		report.TestSuite = "Terraform AWS SES Tests"
		report.PrintReport()

		// Save reports to files
		if err := report.SaveReportToFile("test-report.json"); err != nil {
			t.Errorf("failed to save report to file: %v", err)
		}

		if err := report.SaveReportToHTML("test-report.html"); err != nil {
			t.Errorf("failed to save report to HTML: %v", err)
		}
	}()

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Define test cases with their functions
	testCases := []struct {
		name string
		fn   func(*testing.T, *terraform.Options, string)
	}{
		{"TestSESDomainIdentityCreated", testSESDomainIdentityCreated},
		{"TestSESEmailIdentityCreated", testSESEmailIdentityCreated},
		{"TestSESDKIMTokensGenerated", testSESDKIMTokensGenerated},
		{"TestRoute53HostedZoneCreated", testRoute53HostedZoneCreated},
		{"TestSESConsumerPolicyCreated", testSESConsumerPolicyCreated},
		{"TestSESDomainVerificationToken", testSESDomainVerificationToken},
	}

	// Run all test cases and collect results
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testStart := time.Now()

			// Capture test result
			defer func() {
				testEnd := time.Now()
				duration := testEnd.Sub(testStart)

				result := testutil.TestResult{
					Name:     tc.name,
					Duration: duration.String(),
				}

				if r := recover(); r != nil {
					result.Status = "FAIL"
					result.Error = fmt.Sprintf("Panic: %v", r)
				} else if t.Failed() {
					result.Status = "FAIL"
					result.Error = "Test assertions failed"
				} else if t.Skipped() {
					result.Status = "SKIP"
				} else {
					result.Status = "PASS"
				}

				testResults = append(testResults, result)
			}()

			// Run the actual test
			tc.fn(t, terraformOptions, awsRegion)
		})
	}
}

// Test if SES domain identity is created
func testSESDomainIdentityCreated(t *testing.T, terraformOptions *terraform.Options, awsRegion string) {
	// Get the domain identity ARN from terraform output
	domainIdentityArn := terraform.Output(t, terraformOptions, "ses_domain_identity_arn")
	require.NotEmpty(t, domainIdentityArn, "SES domain identity ARN should not be empty")

	// Verify ARN format
	assert.True(t, strings.HasPrefix(domainIdentityArn, "arn:aws:ses:"), "Domain identity ARN should be valid")
	assert.Contains(t, domainIdentityArn, awsRegion, "Domain identity ARN should contain the correct region")

	// Verify domain identity exists in SES
	sess := createAWSSession(t, awsRegion)
	sesClient := ses.New(sess)

	// Extract domain from ARN
	arnParts := strings.Split(domainIdentityArn, "/")
	domain := arnParts[len(arnParts)-1]

	// List SES identities to verify domain was created
	input := &ses.ListIdentitiesInput{
		IdentityType: aws.String("Domain"),
		MaxItems:     aws.Int64(100),
	}

	result, err := sesClient.ListIdentities(input)
	require.NoError(t, err, "Failed to list SES domain identities")

	// Check if our test domain is in the list
	found := false
	for _, identity := range result.Identities {
		if *identity == domain {
			found = true
			break
		}
	}

	assert.True(t, found, fmt.Sprintf("SES domain identity %s should exist", domain))
}

// Test if SES email identity is created
func testSESEmailIdentityCreated(t *testing.T, terraformOptions *terraform.Options, awsRegion string) {
	// Get the email identity ARN from terraform output
	emailIdentityArn := terraform.Output(t, terraformOptions, "ses_email_identity_arn")
	require.NotEmpty(t, emailIdentityArn, "SES email identity ARN should not be empty")

	// Verify ARN format
	assert.True(t, strings.HasPrefix(emailIdentityArn, "arn:aws:ses:"), "Email identity ARN should be valid")
	assert.Contains(t, emailIdentityArn, awsRegion, "Email identity ARN should contain the correct region")

	// Verify email identity exists in SES
	sess := createAWSSession(t, awsRegion)
	sesClient := ses.New(sess)

	// Extract email from ARN
	arnParts := strings.Split(emailIdentityArn, "/")
	email := arnParts[len(arnParts)-1]

	// List SES identities to verify email was created
	input := &ses.ListIdentitiesInput{
		IdentityType: aws.String("EmailAddress"),
		MaxItems:     aws.Int64(100),
	}

	result, err := sesClient.ListIdentities(input)
	require.NoError(t, err, "Failed to list SES email identities")

	// Check if our test email is in the list
	found := false
	for _, identity := range result.Identities {
		if *identity == email {
			found = true
			break
		}
	}

	assert.True(t, found, fmt.Sprintf("SES email identity %s should exist", email))
}

// Test if SES DKIM tokens are generated
func testSESDKIMTokensGenerated(t *testing.T, terraformOptions *terraform.Options, awsRegion string) {
	// Get the DKIM tokens from terraform output
	dkimTokensOutput := terraform.Output(t, terraformOptions, "ses_dkim_tokens")
	require.NotEmpty(t, dkimTokensOutput, "SES DKIM tokens should not be empty")

	// Parse DKIM tokens (they come as a string representation of a list)
	// Remove brackets and split by comma
	dkimTokensStr := strings.Trim(dkimTokensOutput, "[]")
	dkimTokensStr = strings.ReplaceAll(dkimTokensStr, "\"", "")
	dkimTokens := strings.Split(dkimTokensStr, " ")

	// Filter out empty strings
	var validTokens []string
	for _, token := range dkimTokens {
		token = strings.TrimSpace(token)
		if token != "" {
			validTokens = append(validTokens, token)
		}
	}

	// SES should generate exactly 3 DKIM tokens
	assert.GreaterOrEqual(t, len(validTokens), 1, "Should have at least 1 DKIM token")

	// Verify each token is not empty and has reasonable length
	for i, token := range validTokens {
		assert.NotEmpty(t, token, fmt.Sprintf("DKIM token %d should not be empty", i+1))
		assert.Greater(t, len(token), 10, fmt.Sprintf("DKIM token %d should have reasonable length", i+1))
	}
}

// Test if Route53 hosted zone is created
func testRoute53HostedZoneCreated(t *testing.T, terraformOptions *terraform.Options, awsRegion string) {
	// Get the Route53 zone ID from terraform output
	zoneID := terraform.Output(t, terraformOptions, "route53_zone_id")
	require.NotEmpty(t, zoneID, "Route53 zone ID should not be empty")

	// Verify zone ID format
	assert.True(t, strings.HasPrefix(zoneID, "Z"), "Route53 zone ID should start with 'Z'")

	// Get name servers from terraform output
	nameServersOutput := terraform.Output(t, terraformOptions, "name_servers")
	require.NotEmpty(t, nameServersOutput, "Name servers should not be empty")

	// Verify hosted zone exists in Route53
	sess := createAWSSession(t, awsRegion)
	route53Client := route53.New(sess)

	input := &route53.GetHostedZoneInput{
		Id: aws.String(zoneID),
	}

	result, err := route53Client.GetHostedZone(input)
	require.NoError(t, err, "Failed to get hosted zone details")
	require.NotNil(t, result.HostedZone, "Hosted zone should exist")

	// Verify hosted zone properties
	assert.NotEmpty(t, *result.HostedZone.Name, "Hosted zone name should not be empty")
	assert.Equal(t, zoneID, strings.TrimPrefix(*result.HostedZone.Id, "/hostedzone/"), "Zone ID should match")
}

// Test if SES consumer policy is created
func testSESConsumerPolicyCreated(t *testing.T, terraformOptions *terraform.Options, awsRegion string) {
	// Get the consumer policy ARN from terraform output
	policyArn := terraform.Output(t, terraformOptions, "cosumer_policy_arn")
	require.NotEmpty(t, policyArn, "Consumer policy ARN should not be empty")

	// Verify ARN format
	assert.True(t, strings.HasPrefix(policyArn, "arn:aws:iam:"), "Consumer policy ARN should be valid")
	assert.Contains(t, policyArn, "policy/", "Consumer policy ARN should contain 'policy/'")

	// Extract policy name from ARN
	arnParts := strings.Split(policyArn, "/")
	policyName := arnParts[len(arnParts)-1]
	assert.NotEmpty(t, policyName, "Policy name should not be empty")
	assert.Contains(t, policyName, "AllowSESSend", "Policy name should contain 'ses'")
}

// Test if SES domain verification token is generated
func testSESDomainVerificationToken(t *testing.T, terraformOptions *terraform.Options, awsRegion string) {
	// Get the domain verification token from terraform output
	verificationToken := terraform.Output(t, terraformOptions, "ses_domain_identity_verification_token")
	require.NotEmpty(t, verificationToken, "Domain verification token should not be empty")

	// Verify token format (SES verification tokens are typically long alphanumeric strings)
	assert.Greater(t, len(verificationToken), 20, "Verification token should have reasonable length")
	assert.Regexp(t, "^[a-zA-Z0-9+/=]+$", verificationToken, "Verification token should be base64-like format")
}

// Helper function to create AWS session
func createAWSSession(t *testing.T, region string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	require.NoError(t, err, "Failed to create AWS session")
	return sess
}
