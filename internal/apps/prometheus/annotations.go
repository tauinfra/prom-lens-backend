package prometheus

import (
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

var reservedAnnotationKeys = map[string]struct{}{
	"summary":     {},
	"description": {},
}

// ValidateExtraAnnotations 校验扩展 annotations 必须为 JSON 对象，且不得包含保留 key。
func ValidateExtraAnnotations(extra datatypes.JSON) error {
	if len(extra) == 0 || string(extra) == "null" {
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(extra, &m); err != nil {
		return fmt.Errorf("extraAnnotations must be a JSON object")
	}
	for k := range m {
		if _, reserved := reservedAnnotationKeys[k]; reserved {
			return fmt.Errorf("extraAnnotations must not contain reserved key: %s", k)
		}
	}
	return nil
}

// NormalizeExtraAnnotations 将空值规范为 {}。
func NormalizeExtraAnnotations(extra datatypes.JSON) datatypes.JSON {
	if len(extra) == 0 || string(extra) == "null" {
		return datatypes.JSON([]byte("{}"))
	}
	return extra
}

// MergeRuleAnnotations 合并固定字段与扩展 annotations，供 Prometheus YAML 使用。
func MergeRuleAnnotations(summary, description string, extra datatypes.JSON) (datatypes.JSON, error) {
	merged := map[string]string{
		"summary":     summary,
		"description": description,
	}
	if len(extra) > 0 && string(extra) != "null" {
		var extraMap map[string]string
		if err := json.Unmarshal(extra, &extraMap); err != nil {
			return nil, fmt.Errorf("invalid extraAnnotations: %w", err)
		}
		for k, v := range extraMap {
			if _, reserved := reservedAnnotationKeys[k]; reserved {
				continue
			}
			merged[k] = v
		}
	}
	raw, err := json.Marshal(merged)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(raw), nil
}

// SplitRuleAnnotations 将 Prometheus annotations 拆为 summary、description 与扩展字段。
func SplitRuleAnnotations(merged datatypes.JSON) (summary, description string, extra datatypes.JSON, err error) {
	extra = datatypes.JSON([]byte("{}"))
	if len(merged) == 0 || string(merged) == "null" {
		return "", "", extra, nil
	}
	var m map[string]string
	if err = json.Unmarshal(merged, &m); err != nil {
		return "", "", nil, fmt.Errorf("annotations must be a JSON object: %w", err)
	}
	extraMap := make(map[string]string, len(m))
	for k, v := range m {
		switch k {
		case "summary":
			summary = v
		case "description":
			description = v
		default:
			extraMap[k] = v
		}
	}
	raw, err := json.Marshal(extraMap)
	if err != nil {
		return "", "", nil, err
	}
	return summary, description, datatypes.JSON(raw), nil
}
