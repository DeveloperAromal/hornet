package fraud

import "time"

type RequestFeatures struct {
	ClientIP     string
	Method       string
	Path         string
	UserAgent    string
	Identity     string
	StatusCode   int
	RequestBytes int64
	Timestamp    time.Time
}

type Finding struct {
	RuleID      string
	Description string
	Score       int
}

type Decision struct {
	RiskScore int
	Findings  []Finding
	Action    string
}
