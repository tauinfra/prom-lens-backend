package service

import (
	"fmt"
	"strings"

	"prom-lens-backend/internal/apps/alerting/model"
	"prom-lens-backend/internal/apps/alerting/request"
)

var allowedMatcherOperators = map[string]struct{}{
	"=": {}, "!=": {}, "=~": {}, "!~": {},
}

func buildRouteFromInput(in *request.RouteInput) (*model.Route, []model.RouteMatcher, error) {
	if in == nil || len(in.Matchers) == 0 {
		return nil, nil, nil
	}
	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}
	priority := 100
	if in.Priority != nil {
		priority = *in.Priority
	}
	cont := false
	if in.Continue != nil {
		cont = *in.Continue
	}
	route := &model.Route{
		Enabled:       &enabled,
		Priority:      priority,
		RouteContinue: &cont,
	}
	matchers := make([]model.RouteMatcher, 0, len(in.Matchers))
	for _, m := range in.Matchers {
		label := strings.TrimSpace(m.Label)
		value := strings.TrimSpace(m.Value)
		if label == "" || value == "" {
			return nil, nil, fmt.Errorf("route matcher label and value are required")
		}
		op := strings.TrimSpace(m.Operator)
		if op == "" {
			op = "="
		}
		if _, ok := allowedMatcherOperators[op]; !ok {
			return nil, nil, fmt.Errorf("invalid matcher operator %q", op)
		}
		matchers = append(matchers, model.RouteMatcher{
			Label:    label,
			Operator: op,
			Value:    value,
		})
	}
	return route, matchers, nil
}
