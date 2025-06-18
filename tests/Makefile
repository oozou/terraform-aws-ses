.PHONY: test generate-report clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  test           - Run all tests"
	@echo "  generate-report - Generate test reports (JSON and HTML)"
	@echo "  clean          - Clean up test artifacts"
	@echo "  help           - Show this help message"

# Run tests with report generation
test:
	@echo "Running Terraform AWS SES tests..."
	go test -v -timeout 45m -args -report=true -report-file=test-report.json -html-file=test-report.html

# Generate test report (used by CI/CD)
generate-report:
	@echo "Generating test reports..."
	@if [ -f test-report.json ]; then \
		echo "Test report already exists: test-report.json"; \
		echo "HTML report already exists: test-report.html"; \
	else \
		echo "No test report found. Running tests to generate report..."; \
		$(MAKE) test; \
	fi

# Clean up test artifacts
clean:
	@echo "Cleaning up test artifacts..."
	rm -f test-report.json
	rm -f test-report.html
	rm -f test-results.json
	rm -f test-summary.md
	@echo "Clean up completed."

# Run tests in verbose mode
test-verbose:
	@echo "Running tests in verbose mode..."
	go test -v -timeout 45m -args -report=true -report-file=test-report.json -html-file=test-report.html

# Run specific test
test-specific:
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-specific TEST=TestName"; \
		exit 1; \
	fi
	@echo "Running specific test: $(TEST)"
	go test -v -timeout 45m -run $(TEST) -args -report=true -report-file=test-report.json -html-file=test-report.html
