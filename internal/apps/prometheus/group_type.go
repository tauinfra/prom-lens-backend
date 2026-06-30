package prometheus

const (
	GroupTypeAlertingRules   = "ALERTING RULES"
	GroupTypeAlertingRecords = "ALERTING RECORDS"
)

func IsAlertingRules(groupType string) bool {
	return groupType == GroupTypeAlertingRules
}

func IsAlertingRecords(groupType string) bool {
	return groupType == GroupTypeAlertingRecords
}
