package executor

import "fmt"

func formatAlertmanagerMatcher(label, operator, value string) string {
	switch operator {
	case "=", "!=", "=~", "!~":
		return fmt.Sprintf(`%s%s"%s"`, label, operator, value)
	default:
		return fmt.Sprintf(`%s="%s"`, label, value)
	}
}
