package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/yudai/pp"
)

const (
	Endpoint    = "http://localhost:9090/execute"
	ContentType = "application/json"

	TotalIterationCount       = 4
	CycleStartingRequestCount = 20
	MaxIncreaseCoefficient    = 4
)

type serverResponse struct {
	BuildBody  string `json:"buildBody"`
	BuildError bool   `json:"buildError"`
	RunBody    string `json:"runBody"`
	RunError   bool   `json:"runError"`
	TestBody   string `json:"testBody"`
	TestError  bool   `json:"testError"`
}

type report struct {
	Garbage       bool
	GarbageReason string
	Cause         string
	Failed        bool
	ExecutionTime time.Duration
}

type requestBody struct {
	Lang    string `json:"language"`
	Content string `json:"content"`
}

var requestBodies = []requestBody{
	{
		Lang:    "python",
		Content: "for i in range(10):\n\tprint('*' * (i + 1))",
	},
	{
		Lang:    "python",
		Content: "for i in range(20):\n\tprint('*' * (i + 1))",
	},
	{
		Lang:    "python",
		Content: "for i in range(30):\n\tprint('*' * (i + 1))",
	},
	{
		Lang:    "python",
		Content: "for i in range(40):\n\tprint('*' * (i + 1))",
	},
}

func main() {
	var iterationReports [][]report

	for i := 0; i < TotalIterationCount; i++ {

		randomIncreaseCoefficient := rand.Intn(MaxIncreaseCoefficient) + 1
		totalRequestCountForCycle := CycleStartingRequestCount * randomIncreaseCoefficient

		var wg sync.WaitGroup
		wg.Add(totalRequestCountForCycle)

		reports := make([]report, totalRequestCountForCycle)
		log.Printf("Iteration %d starts with %d concurrent requests\n", i, totalRequestCountForCycle)
		for j := 0; j < totalRequestCountForCycle; j++ {
			go hitAndRun(&wg, j, reports)
		}
		wg.Wait()
		iterationReports = append(iterationReports, reports)

		log.Printf("Iteration %d finished..\nWaiting for a second..\n", i)
		time.Sleep(time.Second)
	}

	summaries := calculateReportsSummaries(iterationReports)
	pp.Println(summaries)
}

func hitAndRun(wg *sync.WaitGroup, idx int, reports []report) {
	singleReport := hit(idx % len(requestBodies))
	reports[idx] = singleReport
	wg.Done()
}

type reportsSummary struct {
	MaxExecutionTimeMs int64
	AvgExecutionTimeMs int64
	MinExecutionTimeMs int64

	TotalCount   int64
	GarbageCount int64
	FailCount    int64
	SuccessCount int64

	SuccessRate    float64
	FailRate       float64
	GarbageRate    float64
	GarbageReasons []string
}

func calculateReportsSummaries(iterationReports [][]report) []reportsSummary {
	var summaries []reportsSummary

	for _, reports := range iterationReports {
		summary := reportsSummary{
			MinExecutionTimeMs: math.MaxInt64,
		}

		for idx := range reports {
			report := reports[idx]
			if report.Failed {
				summary.FailCount++
			} else if report.Garbage {
				summary.GarbageCount++
				summary.GarbageReasons = append(summary.GarbageReasons, report.GarbageReason)
			} else {
				summary.SuccessCount++
			}

			summary.MaxExecutionTimeMs = max(summary.MaxExecutionTimeMs, report.ExecutionTime.Milliseconds())
			summary.MinExecutionTimeMs = min(summary.MinExecutionTimeMs, report.ExecutionTime.Milliseconds())
			summary.AvgExecutionTimeMs += report.ExecutionTime.Milliseconds()
		}

		summary.TotalCount = summary.FailCount + summary.SuccessCount + summary.GarbageCount
		summary.SuccessRate = float64(summary.SuccessCount) / float64(summary.TotalCount) * 100
		summary.FailRate = float64(summary.FailCount) / float64(summary.TotalCount) * 100
		summary.GarbageRate = float64(summary.GarbageCount) / float64(summary.TotalCount) * 100
		summary.AvgExecutionTimeMs = summary.AvgExecutionTimeMs / summary.TotalCount

		summaries = append(summaries, summary)
	}

	return summaries
}

func min(t1, t2 int64) int64 {
	if t1 < t2 {
		return t1
	}
	return t2
}

func max(t1, t2 int64) int64 {
	if t1 > t2 {
		return t1
	}
	return t2
}

func hit(idx int) report {
	var report report
	randomRequestBody := requestBodies[idx]
	reqBytes, _ := json.Marshal(randomRequestBody)

	startTime := time.Now()
	res, err := http.DefaultClient.Post(Endpoint, ContentType, bytes.NewBuffer(reqBytes))
	report.ExecutionTime = time.Since(startTime)
	if err != nil {
		report.Garbage = true
		report.GarbageReason = fmt.Sprintf("%v", err)
		return report
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		report.Failed = true
		report.Cause = err.Error()
		return report
	}

	var response serverResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		report.Cause = err.Error()
		report.Failed = true
		return report
	}

	isServerResponseOK := res.StatusCode >= 200 && res.StatusCode < 400
	if !isServerResponseOK {
		report.Cause = response.RunBody
		report.Failed = true
	}

	return report
}
