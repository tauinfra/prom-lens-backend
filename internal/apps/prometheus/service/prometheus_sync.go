package service

import (
	"context"
	"errors"
	"fmt"

	prom "prom-lens-backend/internal/apps/prometheus"
	"prom-lens-backend/internal/apps/prometheus/dto"
	"prom-lens-backend/internal/apps/prometheus/executor"
	"prom-lens-backend/internal/apps/prometheus/model"
	"prom-lens-backend/internal/apps/prometheus/repository"
	"prom-lens-backend/internal/core/logger"

	"gorm.io/gorm"
)

// SyncManager ConfigMap 与数据库同步（导入）。
type SyncManager interface {
	ImportRulesFromConfigMap(ctx context.Context) (dto.ImportRulesResult, error)
}

type syncManager struct {
	group  repository.GroupRepository
	rule   repository.RuleRepository
	record repository.RecordRepository
}

func NewSyncManager(
	group repository.GroupRepository,
	rule repository.RuleRepository,
	record repository.RecordRepository,
) SyncManager {
	return &syncManager{group: group, rule: rule, record: record}
}

func (s *syncManager) ImportRulesFromConfigMap(ctx context.Context) (dto.ImportRulesResult, error) {
	result := dto.ImportRulesResult{}

	logger.Info("[prom-import] start full configmap rules import")

	entries, err := executor.ListRuleConfigMapEntries(ctx)
	if err != nil {
		logPromImportFailed("list_configmap", "", err)
		return result, err
	}
	result.FilesTotal = len(entries)

	for key, content := range entries {
		parsedGroups, err := executor.ParseRuleGroupsFromYAML(key, []byte(content))
		if err != nil {
			result.Errors = append(result.Errors, dto.ImportErrorDTO{
				ConfigMapKey: key,
				Message:      err.Error(),
			})
			continue
		}
		for _, parsed := range parsedGroups {
			s.importParsedGroup(ctx, parsed, &result)
		}
	}

	logger.Infof(
		"[prom-import] done files=%d groupsCreated=%d groupsSkipped=%d rulesCreated=%d rulesSkipped=%d recordsCreated=%d recordsSkipped=%d errors=%d",
		result.FilesTotal,
		result.GroupsCreated,
		result.GroupsSkipped,
		result.RulesCreated,
		result.RulesSkipped,
		result.RecordsCreated,
		result.RecordsSkipped,
		len(result.Errors),
	)
	return result, nil
}

func (s *syncManager) importParsedGroup(ctx context.Context, parsed *executor.ImportRuleGroup, result *dto.ImportRulesResult) {
	group, err := s.ensureGroup(ctx, parsed, result)
	if err != nil {
		return
	}

	for _, entry := range parsed.Rules {
		switch entry.Kind {
		case executor.RuleEntryAlert:
			s.importAlertRule(ctx, group.ID, parsed.ConfigMapKey, parsed.GroupName, entry, result)
		case executor.RuleEntryRecord:
			s.importRecord(ctx, group.ID, parsed.ConfigMapKey, parsed.GroupName, entry, result)
		default:
			result.Errors = append(result.Errors, dto.ImportErrorDTO{
				ConfigMapKey: parsed.ConfigMapKey,
				GroupName:    parsed.GroupName,
				ItemName:     entry.Name,
				Message:      fmt.Sprintf("unsupported rule kind: %s", entry.Kind),
			})
		}
	}
}

func (s *syncManager) ensureGroup(ctx context.Context, parsed *executor.ImportRuleGroup, result *dto.ImportRulesResult) (model.Group, error) {
	existing, err := s.group.GetByName(ctx, parsed.GroupName)
	if err == nil {
		if existing.Type != parsed.GroupType {
			result.Errors = append(result.Errors, dto.ImportErrorDTO{
				ConfigMapKey: parsed.ConfigMapKey,
				GroupName:    parsed.GroupName,
				Message: fmt.Sprintf(
					"group already exists with type %q, configmap file type is %q",
					existing.Type, parsed.GroupType,
				),
			})
			return model.Group{}, fmt.Errorf("group type mismatch")
		}
		result.GroupsSkipped++
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: parsed.ConfigMapKey,
			GroupName:    parsed.GroupName,
			Message:      err.Error(),
		})
		return model.Group{}, err
	}

	data := &model.Group{
		Name:        parsed.GroupName,
		Type:        parsed.GroupType,
		Description: fmt.Sprintf("imported from configmap %s", parsed.ConfigMapKey),
	}
	if err := s.group.Create(ctx, data); err != nil {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: parsed.ConfigMapKey,
			GroupName:    parsed.GroupName,
			Message:      err.Error(),
		})
		return model.Group{}, err
	}
	result.GroupsCreated++
	logger.Infof("[prom-import] created group name=%q type=%s", parsed.GroupName, parsed.GroupType)
	return *data, nil
}

func (s *syncManager) importAlertRule(
	ctx context.Context,
	groupID int,
	configMapKey, groupName string,
	entry executor.ImportRuleEntry,
	result *dto.ImportRulesResult,
) {
	if _, err := s.rule.GetByGroupAndName(ctx, groupID, entry.Name); err == nil {
		result.RulesSkipped++
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			ItemName:     entry.Name,
			Message:      err.Error(),
		})
		return
	}

	summary, description, extra, err := prom.SplitRuleAnnotations(entry.Annotations)
	if err != nil {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			ItemName:     entry.Name,
			Message:      err.Error(),
		})
		return
	}
	if summary == "" {
		summary = entry.Name
	}
	if description == "" {
		description = "-"
	}
	if err := prom.ValidateExtraAnnotations(extra); err != nil {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			ItemName:     entry.Name,
			Message:      err.Error(),
		})
		return
	}

	enabled := true
	data := &model.Rule{
		Name:             entry.Name,
		GroupID:          groupID,
		Summary:          summary,
		Description:      description,
		Expr:             entry.Expr,
		For:              entry.For,
		Labels:           entry.Labels,
		ExtraAnnotations: prom.NormalizeExtraAnnotations(extra),
		Status:           &enabled,
	}
	if err := s.rule.Create(ctx, data); err != nil {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			ItemName:     entry.Name,
			Message:      err.Error(),
		})
		return
	}
	result.RulesCreated++
}

func (s *syncManager) importRecord(
	ctx context.Context,
	groupID int,
	configMapKey, groupName string,
	entry executor.ImportRuleEntry,
	result *dto.ImportRulesResult,
) {
	if _, err := s.record.GetByGroupAndName(ctx, groupID, entry.Name); err == nil {
		result.RecordsSkipped++
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			ItemName:     entry.Name,
			Message:      err.Error(),
		})
		return
	}

	data := &model.Record{
		Name:    entry.Name,
		GroupID: groupID,
		Expr:    entry.Expr,
	}
	if err := s.record.Create(ctx, data); err != nil {
		result.Errors = append(result.Errors, dto.ImportErrorDTO{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			ItemName:     entry.Name,
			Message:      err.Error(),
		})
		return
	}
	result.RecordsCreated++
}

func logPromImportFailed(step, key string, err error) {
	logger.Errorf("[prom-import] failed step=%s configMapKey=%q err=%v", step, key, err)
}
