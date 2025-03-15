package report

import (
	"fmt"
	"sort"
	"time"

	"stress-tester/internal/dto"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ReportRed takes a map[string]*dto.ResultRed and prints a report of the
// results from the test run. The report is sorted by time and formatted in a
// human-readable format. The columns are:
//
// - Rate: The number of requests per second
// - Error: The number of requests that had an error
// - Avg Time: The average time taken for the requests, excluding network errors
// - Min Time: The minimum time taken for the requests, excluding network errors
// - Max Time: The maximum time taken for the requests, excluding network errors
// - Net Error: The number of requests that had a network error
func ReportRed(result map[string]*dto.ResultRed) {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	p := message.NewPrinter(language.English)
	fmt.Printf("%10s\t%10s\t%10s\t%10s\t%10s\t%10s\n", "Rate", "Error", "Avg Time", "Min Time", "Max Time", "Net Error")
	for _, v := range keys {
		fmt.Printf("%10s\t%10s\t%10v\t%10v\t%10v\t%10s\n", p.Sprintf("%d", result[v].NumRequestPerSecond), p.Sprintf("%d", result[v].NumRequestWithErrorPerSecond), result[v].AverageDuration/time.Duration(result[v].NumRequestPerSecond-result[v].NumNetworkErrorPerSecond+1), result[v].MinDuration, result[v].MaxDuration, p.Sprintf("%d", result[v].NumNetworkErrorPerSecond))
	}
}

// ReportError takes a map[int]*dto.ResultError and prints a report of the number of times each status code was encountered
// during the test run. The report is sorted by status code and formatted in a human-readable format.
func ReportError(errors map[int]*dto.ResultError) {
	keys := make([]int, 0, len(errors))
	for k := range errors {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	p := message.NewPrinter(language.English)
	fmt.Printf("\n%-7s\t%10s\n", "Status", "# Responses")
	for _, v := range keys {
		fmt.Printf("%-7s\t%10s\n", p.Sprintf("%d", v), p.Sprintf("%d", errors[v].NumRequestWithErrorPerSecond))
	}
}

// ReportPercentiles prints the percentiles for a given dto.Percentiles
// to the console, in a human-readable format.
func ReportPercentiles(perc dto.Percentiles) {
	fmt.Printf("\n%-10s\t%10s\n", "Percentile", "Duration")
	fmt.Printf("%-10s\t%10v\n", "P10", perc.P10)
	fmt.Printf("%-10s\t%10v\n", "P25", perc.P25)
	fmt.Printf("%-10s\t%10v\n", "P50", perc.P50)
	fmt.Printf("%-10s\t%10v\n", "P75", perc.P75)
	fmt.Printf("%-10s\t%10v\n", "P90", perc.P90)
	fmt.Printf("%-10s\t%10v\n", "P99", perc.P99)
}
