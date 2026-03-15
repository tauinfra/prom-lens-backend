package dto

import "valyria-backend/internal/apps/prometheus/model"

func ToGroupDTO(item model.Group) GroupDTO {
	return GroupDTO{
		ID:          item.ID,
		Name:        item.Name,
		Type:        item.Type,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ToRecordDTO(item model.Record) RecordDTO {
	return RecordDTO{
		ID:        item.ID,
		Name:      item.Name,
		GroupID:   item.GroupID,
		Expr:      item.Expr,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func ToRuleDTO(item model.Rule) RuleDTO {
	return RuleDTO{
		ID:          item.ID,
		Name:        item.Name,
		GroupID:     item.GroupID,
		Summary:     item.Summary,
		Description: item.Description,
		Expr:        item.Expr,
		For:         item.For,
		Labels:      item.Labels,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func ToTargetDTO(item model.Target) TargetDTO {
	return TargetDTO{
		ID:        item.ID,
		GroupID:   item.GroupID,
		IPAddress: item.IPAddress,
		Port:      item.Port,
		Labels:    item.Labels,
		Enabled:   item.Enabled,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func ToTargetGroupDTO(item model.TargetGroup) TargetGroupDTO {
	return TargetGroupDTO{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Labels:      item.Labels,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
