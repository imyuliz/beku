package beku

const (
	qosKey     = "qos"
	autoQosKey = "autoQos"
)

// qos rank,the higher the number, the higher the level
const (
	BestEffortRank = iota
	BurstableRank
	GuaranteedRank
)

var (
	qosRanks = map[string]int{
		"BestEffort": BestEffortRank,
		"Burstable":  BurstableRank,
		"Guaranteed": GuaranteedRank,
	}
)
