package repository

import "knative.dev/pkg/apis"

func tektonSucceededAndReason(conditions []apis.Condition) (string, string) {
	for _, condition := range conditions {
		if condition.Type == apis.ConditionSucceeded {
			return string(condition.Status), condition.Reason
		}
	}
	return "Unknown", ""
}
