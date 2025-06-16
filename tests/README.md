# Terraform AWS SES Module Tests

This directory contains automated tests for the terraform-aws-ses module using Go and Terratest with comprehensive reporting capabilities.

## Test Files

- `terraform_aws_ses_test.go` - Consolidated integration tests with built-in reporting functionality
- `go.mod` - Go module dependencies

## Test Structure

The test suite has been consolidated into a single comprehensive test file that includes:
- **Test execution logic** - Core test functions for email and domain verification
- **Test reporting** - Built-in HTML, JSON, and Markdown report generation
- **GitHub integration** - Automatic PR comment generation with test results

## Running Tests

### Prerequisites

1. Go 1.21+ installed
2. Terraform installed
3. AWS credentials configured

### Basic Test Execution

Run tests without reporting:

```bash
cd tests
export AWS_REGION=ap-southeast-1
go test -v -timeout 30m
```

### Test Execution with Reporting

Run tests with comprehensive reporting:

```bash
cd tests
export AWS_REGION=ap-southeast-1
go test -v -timeout 30m -report -report-file=test-report.json -html-file=test-report.html
```

This will generate:
- `test-report.json` - Detailed JSON test report
- `test-report.html` - Interactive HTML test report
- `test-summary.md` - GitHub-friendly markdown summary
- `test-results.json` - Simplified results for CI/CD integration

**Note:** Integration tests will create and destroy AWS resources, which may incur costs.

## Test Configuration

The tests use the examples with variable overrides:

- **Email Verification**: Tests pass `ses_email` variable with dynamic test email
- **Domain Verification**: Tests pass `ses_domain` and `route53_zone_name` variables with dynamic test domains
- **Route53 Integration**: Domain tests create actual hosted zones for verification

## What the Tests Validate

### SES Email Identity Tests
- ✅ Email identity resource creation
- ✅ Email format validation
- ✅ SES identity registration (integration test)
- ✅ Verification attributes validation

### SES Domain Identity Tests  
- ✅ Domain identity resource creation
- ✅ DKIM configuration and enablement
- ✅ Domain verification setup (integration test)
- ✅ DKIM attributes validation

### DMARC Tests
- ✅ DMARC Route53 record creation
- ✅ TXT record type validation
- ✅ DMARC policy format validation (v=DMARC1)
- ✅ DNS record deployment (integration test)

## Test Reporting Features

### Console Output
- Real-time test progress with emojis and formatting
- Detailed statistics (pass rate, duration, etc.)
- Error details for failed tests
- Summary with actionable insights

### HTML Report
- Interactive web-based test report
- Visual progress bars and statistics
- Color-coded test results
- Responsive design for mobile viewing

### JSON Report
- Machine-readable test results
- Detailed test metadata and timing
- Error information and stack traces
- Suitable for CI/CD integration

### GitHub Integration
- Automatic PR comments with test results
- Markdown-formatted summaries
- Status badges (pass/fail)
- Failed test details with error messages

## CI/CD Integration

Tests are automatically run on pull requests via GitHub Actions workflow with the following features:

### Automated Failure Handling
- **Validation failures**: Terraform format/validation issues trigger "@claude fix build error:" comments
- **Lint failures**: Go linting issues trigger "@claude fix build error:" comments  
- **Test failures**: Failed integration tests trigger "@claude fix build error:" comments
- **Build links**: All failure comments include direct links to the GitHub Actions build logs

### Test Result Comments
- Successful tests post detailed results to PR comments
- Failed tests post error details with "@claude fix build error:" prefix
- Comments are updated (not duplicated) on subsequent runs
- Rich formatting with tables, badges, and error details

### Artifacts
- Test reports are uploaded as GitHub Actions artifacts
- HTML and JSON reports available for download
- Test summaries integrated with GitHub's step summary feature

## Example Test Output

### Console Output
```
🧪 TERRAFORM AWS SES TEST REPORT
================================================================================
📅 Test Suite: Terraform AWS SES Tests
⏰ Start Time: 2024-01-15 10:30:00 UTC
⏰ End Time:   2024-01-15 10:35:30 UTC
⏱️  Duration:   5m30s

📊 TEST STATISTICS
================================================================================
📈 Total Tests:   2
✅ Passed Tests:  2
❌ Failed Tests:  0
⏭️  Skipped Tests: 0
📊 Pass Rate:     100.0%

📋 DETAILED TEST RESULTS
================================================================================
1. TestSESEmailVerification - ✅ PASS (2m15s)
2. TestSESDomainVerification - ✅ PASS (3m15s)

📝 SUMMARY
================================================================================
✅ ALL TESTS PASSED! 2/2 tests successful

🎉 Congratulations! All tests passed successfully!
```

### GitHub PR Comment (Success)
```markdown
## 🧪 Terraform AWS SES Test Results

![Tests Passed](https://img.shields.io/badge/tests-passed-success)

### 📊 Summary

| Metric | Value |
|--------|-------|
| Total Tests | 2 |
| ✅ Passed | 2 |
| ❌ Failed | 0 |
| ⏭️ Skipped | 0 |
| 📊 Pass Rate | 100.0% |
| ⏱️ Duration | 5m30s |

### 📋 Test Details

| Test Name | Status | Duration |
|-----------|--------|----------|
| TestSESEmailVerification | ✅ Pass | 2m15s |
| TestSESDomainVerification | ✅ Pass | 3m15s |

### 📝 Summary

✅ ALL TESTS PASSED! 2/2 tests successful
```

### GitHub PR Comment (Failure)
```markdown
@claude fix build error:

**Build Link:** https://github.com/owner/repo/actions/runs/123456789

## 🧪 Terraform AWS SES Test Results

![Tests Failed](https://img.shields.io/badge/tests-failed-critical)

### ❌ Failed Tests

**TestSESEmailVerification**
```
failed to list SES identities: AccessDenied: User is not authorized to perform: ses:ListIdentities
```

### 📝 Summary

❌ 1/2 tests failed, 1 passed
```

## Troubleshooting

### Common Issues

1. **AWS Credentials**: Ensure AWS credentials are properly configured
2. **Permissions**: Tests require SES, Route53, and IAM permissions
3. **Region**: Make sure AWS_REGION is set to ap-southeast-1 or your preferred region
4. **Timeouts**: Large tests may need increased timeout values

### Debug Mode

For verbose debugging, run tests with additional flags:

```bash
go test -v -timeout 30m -report -args -test.v
