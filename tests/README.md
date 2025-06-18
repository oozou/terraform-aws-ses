# Terraform AWS SES Module Tests

This directory contains comprehensive tests for the Terraform AWS SES module using the `terraform-test-util` framework.

## Overview

The tests validate the functionality of the AWS SES module by:
- Creating SES domain and email identities
- Verifying DKIM token generation
- Testing Route53 hosted zone creation
- Validating IAM consumer policy creation
- Checking domain verification tokens

## Test Framework

The tests use the [terraform-test-util](https://github.com/oozou/terraform-test-util) framework which provides:
- Comprehensive test reporting with JSON and HTML outputs
- Beautiful console output with emojis and clear formatting
- Pass rate calculation and test statistics
- GitHub-friendly summary generation

## Prerequisites

- Go 1.21 or later
- Terraform 1.6.0 or later
- AWS credentials configured
- Access to AWS SES, Route53, and IAM services

## Running Tests

### Using Make (Recommended)

```bash
# Run all tests with report generation
make test

# Generate test reports (if tests were already run)
make generate-report

# Clean up test artifacts
make clean

# Run tests in verbose mode
make test-verbose

# Run a specific test
make test-specific TEST=TestSESDomainIdentityCreated
```

### Using Go Test Directly

```bash
# Run all tests
go test -v -timeout 45m

# Run tests with report generation
go test -v -timeout 45m -args -report=true -report-file=test-report.json -html-file=test-report.html

# Run a specific test
go test -v -timeout 45m -run TestSESDomainIdentityCreated
```

## Test Structure

The test suite includes the following test cases:

1. **TestSESDomainIdentityCreated** - Verifies SES domain identity creation
2. **TestSESEmailIdentityCreated** - Verifies SES email identity creation
3. **TestSESDKIMTokensGenerated** - Validates DKIM token generation
4. **TestRoute53HostedZoneCreated** - Tests Route53 hosted zone creation
5. **TestSESConsumerPolicyCreated** - Verifies IAM consumer policy creation
6. **TestSESDomainVerificationToken** - Checks domain verification token generation

## Test Configuration

The tests use the `examples/terraform-test` configuration which includes:
- Both domain and email SES configurations
- Route53 hosted zone for domain verification
- IAM consumer policy creation
- Proper tagging and naming conventions

## Test Reports

The framework generates multiple types of reports:

### JSON Report (`test-report.json`)
Contains detailed test results in JSON format for programmatic processing.

### HTML Report (`test-report.html`)
A beautiful, interactive HTML report with:
- Test statistics and pass rates
- Detailed test results with status indicators
- Error details for failed tests
- Progress bars and visual indicators

### Console Output
Formatted console output with:
- Test execution progress
- Detailed statistics
- Pass/fail indicators with emojis
- Summary and recommendations

## CI/CD Integration

The tests are integrated with GitHub Actions workflow that:
- Runs tests automatically on pull requests
- Generates and uploads test reports
- Posts test results as PR comments
- Includes "@claude fix build error:" prefix for failed tests
- Provides direct links to workflow runs and detailed logs

## Environment Variables

The following environment variables can be used to configure the tests:

- `AWS_DEFAULT_REGION` - AWS region for testing (default: ap-southeast-1)
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key
- `AWS_SESSION_TOKEN` - AWS session token (if using temporary credentials)

## Troubleshooting

### Common Issues

1. **AWS Permissions**: Ensure your AWS credentials have permissions for SES, Route53, and IAM
2. **Domain Conflicts**: Tests use unique domain names to avoid conflicts
3. **Timeout Issues**: Tests have a 45-minute timeout; adjust if needed
4. **Resource Cleanup**: The `terraform destroy` is called automatically in defer blocks

### Debug Mode

For debugging failed tests:

```bash
# Run with verbose output
go test -v -timeout 45m -run TestSpecificTest

# Check terraform logs
export TF_LOG=DEBUG
go test -v -timeout 45m
```

## Contributing

When adding new tests:

1. Follow the existing test structure and naming conventions
2. Add proper error handling and cleanup
3. Include meaningful assertions and error messages
4. Update this README with new test descriptions
5. Ensure tests are idempotent and don't interfere with each other

## Dependencies

The tests depend on:

- `github.com/gruntwork-io/terratest` - Terraform testing framework
- `github.com/stretchr/testify` - Test assertions and utilities
- `github.com/aws/aws-sdk-go` - AWS SDK for Go
- `github.com/oozou/terraform-test-util` - Test reporting utilities

## License

This test suite is part of the terraform-aws-ses module and follows the same license terms.
