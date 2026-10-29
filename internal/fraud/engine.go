package fraud

import "sync"

type Rule interface {
	ID() string
	Evaluate(RequestFeatures) *Finding
}

type FraudEngine struct {
	mu    sync.RWMutex
	rules []Rule
}

func NewFraudEngine(rules ...Rule) *FraudEngine {
	return &FraudEngine{}
}

func (fe *FraudEngine) Analyze(req RequestFeatures) Decision {
	fe.mu.RLock()
	defer fe.mu.RUnlock()

	decision := Decision{
		Action: "ALLOW",
	}

	for _, rule := range fe.rules {

		finding := rule.Evaluate(req)

		if finding == nil {
			continue
		}

		decision.Findings = append(decision.Findings, *finding)
		decision.RiskScore += finding.Score
	}

	if decision.RiskScore >= 80 {
		decision.Action = "BLOCK"
	} else if decision.RiskScore >= 40 {
		decision.Action = "MONITOR"
	}

	return decision

}
