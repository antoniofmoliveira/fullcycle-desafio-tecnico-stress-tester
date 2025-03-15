package stats

import (
	"sort"

	"stress-tester/internal/dto"
)

// CalculateRed takes a slice of *dto.Red records and returns a map[string]*dto.ResultRed, where the keys are the different second intervals
// and the values are the respective *dto.ResultRed struct containing the average duration, min duration, max duration, and request per second.
func CalculateRed(recs []*dto.Red) map[string]*dto.ResultRed {
	mapRecs := make(map[string]*dto.ResultRed)
	for _, rec := range recs {
		// s := fmt.Sprintf("%d", rec.SentAt.Minute()*100+rec.SentAt.Second())
		s := "0000"
		if mapRecs[s] == nil {
			mapRecs[s] = &dto.ResultRed{}
		}
		mapRecs[s].NumRequestPerSecond++
		if rec.StatusCode == -1 {
			mapRecs[s].NumNetworkErrorPerSecond++
		}
		if rec.StatusCode != 200 && rec.StatusCode != -1 {
			mapRecs[s].NumRequestWithErrorPerSecond++
		}
		mapRecs[s].AverageDuration = mapRecs[s].AverageDuration + rec.ReceivedAt.Sub(rec.SentAt)
		if rec.ReceivedAt.Sub(rec.SentAt) > mapRecs[s].MaxDuration {
			mapRecs[s].MaxDuration = rec.ReceivedAt.Sub(rec.SentAt)
		}
		if rec.ReceivedAt.Sub(rec.SentAt) < mapRecs[s].MinDuration && mapRecs[s].MinDuration != 0 {
			mapRecs[s].MinDuration = rec.ReceivedAt.Sub(rec.SentAt)
		} else if mapRecs[s].MinDuration == 0 {
			mapRecs[s].MinDuration = rec.ReceivedAt.Sub(rec.SentAt)
		}
	}
	return mapRecs
}

// CalculateErrors takes a slice of *dto.Red records and returns a map[int]*dto.ResultError, where the keys are the different status codes
// and the values are the respective *dto.ResultError struct containing the error type and the number of requests with that error per second.
func CalculateErrors(recs []*dto.Red) map[int]*dto.ResultError {

	mapErrs := make(map[int]*dto.ResultError)
	for _, rec := range recs {
		if mapErrs[rec.StatusCode] == nil {
			mapErrs[rec.StatusCode] = &dto.ResultError{ErrorType: rec.StatusCode, NumRequestWithErrorPerSecond: 0}
		}
		mapErrs[rec.StatusCode].NumRequestWithErrorPerSecond++
	}
	return mapErrs
}

// CalculatePercentile calculates the percentiles (P10, P25, P50, P75, P90, P99)
// for the given slice of dto.Red records based on their Duration field.
// It returns a dto.Percentiles struct containing the calculated percentiles.
// The function assumes that the input slice is non-empty.

func CalculatePercentile(recs []*dto.Red) dto.Percentiles {
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Duration < recs[j].Duration
	})
	l := len(recs)

	p10 := int(float64(l) * 10 / 100)
	p25 := int(float64(l) * 25 / 100)
	p50 := int(float64(l) * 50 / 100)
	p75 := int(float64(l) * 75 / 100)
	p90 := int(float64(l) * 90 / 100)
	p99 := int(float64(l) * 99 / 100)

	per := dto.Percentiles{
		P10: recs[p10].Duration,
		P25: recs[p25].Duration,
		P50: recs[p50].Duration,
		P75: recs[p75].Duration,
		P90: recs[p90].Duration,
		P99: recs[p99].Duration,
	}
	return per

}
