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

// qosNotices set Qos information
var (
	qosNotices = map[string]string{
		"Guaranteed": "Every Container in the Pod must have a memory limit and a memory request, and they must be the same, Every Container in the Pod must have a CPU limit and a CPU request, and they must be the same,more information: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod",
		"Burstable":  "The Pod does not meet the criteria for QoS class Guaranteed and at least one Container in the Pod has a memory or CPU request, more information: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod",
		"BestEffort": "The Containers in the Pod must not have any memory or CPU limits or requests, more information: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod",
	}
)
