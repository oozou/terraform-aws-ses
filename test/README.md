# Terraform AWS SES Module Tests

This directory contains automated tests for the terraform-aws-ses module using Go and Terratest.

## Test Files

- `terraform_aws_ses_test.go` - Integration tests that validate actual AWS resources (requires AWS credentials)
- `go.mod` - Go module dependencies

## Running Tests

### Prerequisites

1. Go 1.21+ installed
2. Terraform installed
3. AWS credentials configured

### Integration Tests

These tests create and validate actual AWS resources using configurable variables:

```bash
cd test
export AWS_REGION=us-east-1
go test -v -run TestTerraformAWSSESModule
```

**Note:** Integration tests will create and destroy AWS resources, which may incur costs.

### Test Configuration

The tests use the examples with variable overrides:

- **Email Verification**: Tests pass `ses_email` variable with dynamic test email
- **Domain Verification**: Tests pass `ses_domain` and `route53_zone_name` variables with dynamic test domains
- **Route53 Integration**: Domain tests create actual hosted zones for verification

## What the Tests Validate

### SES Email Identity Tests
- ✅ Email identity resource creation
- ✅ Email format validation
- ✅ SES identity registration (integration test)

### SES Domain Identity Tests  
- ✅ Domain identity resource creation
- ✅ DKIM configuration
- ✅ Domain verification setup (integration test)

### DMARC Tests
- ✅ DMARC Route53 record creation
- ✅ TXT record type validation
- ✅ DMARC policy format validation
- ✅ DNS record deployment (integration test)

## CI/CD Integration

Tests are automatically run on pull requests via GitHub Actions workflow.