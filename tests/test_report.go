package test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// TestResult represents the result of a single test
type TestResult struct {
	Name     string `json:"name"`
	Status   string `json:"status"` // "PASS", "FAIL", "SKIP"
	Duration string `json:"duration"`
	Error    string `json:"error,omitempty"`
}

// TestReport represents the overall test report
type TestReport struct {
	TestSuite    string       `json:"test_suite"`
	StartTime    time.Time    `json:"start_time"`
	EndTime      time.Time    `json:"end_time"`
	Duration     string       `json:"duration"`
	TotalTests   int          `json:"total_tests"`
	PassedTests  int          `json:"passed_tests"`
	FailedTests  int          `json:"failed_tests"`
	SkippedTests int          `json:"skipped_tests"`
	Results      []TestResult `json:"results"`
	Summary      string       `json:"summary"`
}

// GenerateTestReport creates a comprehensive test report
func GenerateTestReport(results []TestResult, startTime, endTime time.Time) *TestReport {
	report := &TestReport{
		TestSuite: "Terraform AWS Aurora PostgreSQL Tests",
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  endTime.Sub(startTime).String(),
		Results:   results,
	}

	// Calculate statistics
	report.TotalTests = len(results)
	for _, result := range results {
		switch result.Status {
		case "PASS":
			report.PassedTests++
		case "FAIL":
			report.FailedTests++
		case "SKIP":
			report.SkippedTests++
		}
	}

	// Generate summary
	if report.FailedTests == 0 {
		report.Summary = fmt.Sprintf("✅ ALL TESTS PASSED! %d/%d tests successful", 
			report.PassedTests, report.TotalTests)
	} else {
		report.Summary = fmt.Sprintf("❌ %d/%d tests failed, %d passed", 
			report.FailedTests, report.TotalTests, report.PassedTests)
	}

	return report
}

// PrintReport prints a formatted test report to console
func (r *TestReport) PrintReport() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🧪 TERRAFORM AWS AURORA POSTGRESQL TEST REPORT")
	fmt.Println(strings.Repeat("=", 80))
	
	fmt.Printf("📅 Test Suite: %s\n", r.TestSuite)
	fmt.Printf("⏰ Start Time: %s\n", r.StartTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("⏰ End Time:   %s\n", r.EndTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("⏱️  Duration:   %s\n", r.Duration)
	
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("📊 TEST STATISTICS")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("📈 Total Tests:   %d\n", r.TotalTests)
	fmt.Printf("✅ Passed Tests:  %d\n", r.PassedTests)
	fmt.Printf("❌ Failed Tests:  %d\n", r.FailedTests)
	fmt.Printf("⏭️  Skipped Tests: %d\n", r.SkippedTests)
	
	if r.TotalTests > 0 {
		passRate := float64(r.PassedTests) / float64(r.TotalTests) * 100
		fmt.Printf("📊 Pass Rate:     %.1f%%\n", passRate)
	}
	
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("📋 DETAILED TEST RESULTS")
	fmt.Println(strings.Repeat("-", 80))
	
	for i, result := range r.Results {
		status := ""
		switch result.Status {
		case "PASS":
			status = "✅ PASS"
		case "FAIL":
			status = "❌ FAIL"
		case "SKIP":
			status = "⏭️  SKIP"
		}
		
		fmt.Printf("%d. %s - %s (%s)\n", i+1, result.Name, status, result.Duration)
		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}
	}
	
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("📝 SUMMARY")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println(r.Summary)
	
	if r.FailedTests == 0 {
		fmt.Println("\n🎉 Congratulations! All PostgreSQL Aurora tests passed successfully!")
		fmt.Println("✅ Your terraform-aws-aurora module is working correctly for PostgreSQL.")
	} else {
		fmt.Println("\n⚠️  Some tests failed. Please review the errors above and fix the issues.")
		fmt.Println("💡 Check the test logs for more detailed error information.")
	}
	
	fmt.Println("\n" + strings.Repeat("=", 80))
}

// SaveReportToFile saves the test report to a JSON file
func (r *TestReport) SaveReportToFile(filename string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write report file: %v", err)
	}
	
	fmt.Printf("📄 Test report saved to: %s\n", filename)
	return nil
}

// SaveReportToHTML saves the test report to an HTML file
func (r *TestReport) SaveReportToHTML(filename string) error {
	html := r.generateHTML()
	
	err := os.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HTML report file: %v", err)
	}
	
	fmt.Printf("🌐 HTML test report saved to: %s\n", filename)
	return nil
}

// generateHTML creates an HTML representation of the test report
func (r *TestReport) generateHTML() string {
	passRate := float64(r.PassedTests) / float64(r.TotalTests) * 100
	
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Terraform AWS Aurora PostgreSQL Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; border-bottom: 2px solid #007bff; padding-bottom: 20px; margin-bottom: 30px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: #f8f9fa; padding: 20px; border-radius: 8px; text-align: center; border-left: 4px solid #007bff; }
        .stat-number { font-size: 2em; font-weight: bold; color: #007bff; }
        .stat-label { color: #6c757d; margin-top: 5px; }
        .test-results { margin-top: 30px; }
        .test-item { display: flex; justify-content: space-between; align-items: center; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .test-pass { background-color: #d4edda; border-left: 4px solid #28a745; }
        .test-fail { background-color: #f8d7da; border-left: 4px solid #dc3545; }
        .test-skip { background-color: #fff3cd; border-left: 4px solid #ffc107; }
        .status { font-weight: bold; padding: 5px 10px; border-radius: 3px; color: white; }
        .status-pass { background-color: #28a745; }
        .status-fail { background-color: #dc3545; }
        .status-skip { background-color: #ffc107; }
        .summary { background: #e9ecef; padding: 20px; border-radius: 8px; margin-top: 30px; text-align: center; }
        .error-details { color: #dc3545; font-size: 0.9em; margin-top: 5px; }
        .progress-bar { width: 100%%; height: 20px; background-color: #e9ecef; border-radius: 10px; overflow: hidden; margin: 10px 0; }
        .progress-fill { height: 100%%; background-color: #28a745; transition: width 0.3s ease; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 Terraform AWS Aurora PostgreSQL Test Report</h1>
            <p><strong>Test Suite:</strong> %s</p>
            <p><strong>Execution Time:</strong> %s to %s</p>
            <p><strong>Duration:</strong> %s</p>
        </div>
        
        <div class="stats">
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Total Tests</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" style="color: #28a745;">%d</div>
                <div class="stat-label">Passed</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" style="color: #dc3545;">%d</div>
                <div class="stat-label">Failed</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" style="color: #ffc107;">%d</div>
                <div class="stat-label">Skipped</div>
            </div>
        </div>
        
        <div>
            <h3>Pass Rate: %.1f%%</h3>
            <div class="progress-bar">
                <div class="progress-fill" style="width: %.1f%%;"></div>
            </div>
        </div>
        
        <div class="test-results">
            <h2>📋 Test Results</h2>`,
		r.TestSuite,
		r.StartTime.Format("2006-01-02 15:04:05 MST"),
		r.EndTime.Format("2006-01-02 15:04:05 MST"),
		r.Duration,
		r.TotalTests,
		r.PassedTests,
		r.FailedTests,
		r.SkippedTests,
		passRate,
		passRate)
	
	for _, result := range r.Results {
		statusClass := ""
		statusText := ""
		switch result.Status {
		case "PASS":
			statusClass = "test-pass"
			statusText = `<span class="status status-pass">✅ PASS</span>`
		case "FAIL":
			statusClass = "test-fail"
			statusText = `<span class="status status-fail">❌ FAIL</span>`
		case "SKIP":
			statusClass = "test-skip"
			statusText = `<span class="status status-skip">⏭️ SKIP</span>`
		}
		
		html += fmt.Sprintf(`
            <div class="test-item %s">
                <div>
                    <strong>%s</strong>
                    <div style="color: #6c757d; font-size: 0.9em;">Duration: %s</div>
                    %s
                </div>
                %s
            </div>`,
			statusClass,
			result.Name,
			result.Duration,
			func() string {
				if result.Error != "" {
					return fmt.Sprintf(`<div class="error-details">Error: %s</div>`, result.Error)
				}
				return ""
			}(),
			statusText)
	}
	
	html += fmt.Sprintf(`
        </div>
        
        <div class="summary">
            <h2>📝 Summary</h2>
            <p><strong>%s</strong></p>
        </div>
    </div>
</body>
</html>`, r.Summary)
	
	return html
}
