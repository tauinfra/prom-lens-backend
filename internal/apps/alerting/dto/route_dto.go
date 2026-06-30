package dto

import (
	"prom-lens-backend/internal/apps/alerting/model"
)

type RouteMatcherDTO struct {
	Label    string `json:"label"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type RouteDTO struct {
	Enabled  *bool             `json:"enabled"`
	Priority int               `json:"priority"`
	Continue *bool             `json:"continue"`
	Matchers []RouteMatcherDTO `json:"matchers,omitempty"`
}

func ToRouteDTO(route *model.Route) *RouteDTO {
	if route == nil {
		return nil
	}
	out := &RouteDTO{
		Enabled:  route.Enabled,
		Priority: route.Priority,
		Continue: route.RouteContinue,
	}
	for _, m := range route.Matchers {
		out.Matchers = append(out.Matchers, RouteMatcherDTO{
			Label:    m.Label,
			Operator: m.Operator,
			Value:    m.Value,
		})
	}
	return out
}
