package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/datatypes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	RuleEntryAlert  = "alert"
	RuleEntryRecord = "record"
)

// ImportRuleEntry 从 ConfigMap 解析出的单条规则（alert 或 record）。
type ImportRuleEntry struct {
	Kind        string // alert | record
	Name        string
	Expr        string
	For         string
	Labels      datatypes.JSON
	Annotations datatypes.JSON
}

// ImportRuleGroup 从单个 ConfigMap 文件解析出的规则组。
type ImportRuleGroup struct {
	ConfigMapKey string
	GroupName    string
	GroupType    string // ALERTING RULES | ALERTING RECORDS
	Rules        []ImportRuleEntry
}

type importYAMLRule struct {
	Alert       string         `json:"alert"`
	Record      string         `json:"record"`
	Expr        string         `json:"expr"`
	For         string         `json:"for"`
	Labels      datatypes.JSON `json:"labels"`
	Annotations datatypes.JSON `json:"annotations"`
}

type importYAMLGroup struct {
	Name  string           `json:"name"`
	Rules []importYAMLRule `json:"rules"`
}

type importYAMLConfig struct {
	Groups []importYAMLGroup `json:"groups"`
}

// ListRuleConfigMapEntries 列出规则 ConfigMap 中全部 .yml/.yaml 条目。
func ListRuleConfigMapEntries(ctx context.Context) (map[string]string, error) {
	client, err := clientSet()
	if err != nil {
		return nil, err
	}
	namespace, name := getPromRuleConfig()
	cm, err := client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get configmap %s/%s: %w", namespace, name, err)
	}
	if len(cm.Data) == 0 {
		return map[string]string{}, nil
	}
	out := make(map[string]string)
	for key, content := range cm.Data {
		if strings.HasSuffix(key, ".yml") || strings.HasSuffix(key, ".yaml") {
			out[key] = content
		}
	}
	return out, nil
}

// ParseRuleGroupsFromYAML 解析 ConfigMap 文件，支持多 group；组名以 YAML name 为准。
func ParseRuleGroupsFromYAML(configMapKey string, content []byte) ([]*ImportRuleGroup, error) {
	jsonData, err := yaml.YAMLToJSON(content)
	if err != nil {
		return nil, fmt.Errorf("yaml to json: %w", err)
	}
	var cfg importYAMLConfig
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal rule config: %w", err)
	}
	if len(cfg.Groups) == 0 {
		return nil, fmt.Errorf("no groups in %s", configMapKey)
	}

	fileFallbackName, err := groupNameFromConfigMapKey(configMapKey)
	if err != nil {
		return nil, err
	}

	out := make([]*ImportRuleGroup, 0, len(cfg.Groups))
	for i, yamlGroup := range cfg.Groups {
		groupName := strings.TrimSpace(yamlGroup.Name)
		if groupName == "" {
			if len(cfg.Groups) == 1 {
				groupName = fileFallbackName
			} else {
				return nil, fmt.Errorf("%s groups[%d]: missing name", configMapKey, i)
			}
		}

		entries, groupType, err := classifyYAMLRules(configMapKey, yamlGroup.Rules)
		if err != nil {
			return nil, fmt.Errorf("%s groups[%d] name=%q: %w", configMapKey, i, groupName, err)
		}

		out = append(out, &ImportRuleGroup{
			ConfigMapKey: configMapKey,
			GroupName:    groupName,
			GroupType:    groupType,
			Rules:        entries,
		})
	}
	return out, nil
}

// ParseRuleGroupFromYAML 解析仅含单个 group 的文件（兼容旧调用）。
func ParseRuleGroupFromYAML(configMapKey string, content []byte) (*ImportRuleGroup, error) {
	groups, err := ParseRuleGroupsFromYAML(configMapKey, content)
	if err != nil {
		return nil, err
	}
	if len(groups) != 1 {
		return nil, fmt.Errorf("%s contains %d groups, expected 1", configMapKey, len(groups))
	}
	return groups[0], nil
}

func groupNameFromConfigMapKey(key string) (string, error) {
	switch {
	case strings.HasSuffix(key, ".yml"):
		return strings.TrimSuffix(key, ".yml"), nil
	case strings.HasSuffix(key, ".yaml"):
		return strings.TrimSuffix(key, ".yaml"), nil
	default:
		return "", fmt.Errorf("unsupported config key: %s", key)
	}
}

func classifyYAMLRules(configMapKey string, rules []importYAMLRule) ([]ImportRuleEntry, string, error) {
	var (
		entries   []ImportRuleEntry
		groupType string
	)
	for i, r := range rules {
		entry, entryGroupType, err := classifyYAMLRule(r)
		if err != nil {
			return nil, "", fmt.Errorf("%s rules[%d]: %w", configMapKey, i, err)
		}
		if groupType == "" {
			groupType = entryGroupType
		} else if groupType != entryGroupType {
			return nil, "", fmt.Errorf("%s: mixed alert and record rules in one file", configMapKey)
		}
		entries = append(entries, entry)
	}
	if groupType == "" {
		groupType = AlertingRules
	}
	return entries, groupType, nil
}

func classifyYAMLRule(r importYAMLRule) (ImportRuleEntry, string, error) {
	hasAlert := strings.TrimSpace(r.Alert) != ""
	hasRecord := strings.TrimSpace(r.Record) != ""
	switch {
	case hasAlert && hasRecord:
		return ImportRuleEntry{}, "", fmt.Errorf("rule has both alert and record")
	case !hasAlert && !hasRecord:
		return ImportRuleEntry{}, "", fmt.Errorf("rule missing alert and record")
	case hasAlert:
		return ImportRuleEntry{
			Kind:        RuleEntryAlert,
			Name:        r.Alert,
			Expr:        r.Expr,
			For:         defaultFor(r.For),
			Labels:      normalizeLabelsJSON(r.Labels),
			Annotations: r.Annotations,
		}, AlertingRules, nil
	default:
		return ImportRuleEntry{
			Kind:   RuleEntryRecord,
			Name:   r.Record,
			Expr:   r.Expr,
			Labels: normalizeLabelsJSON(r.Labels),
		}, AlertingRecords, nil
	}
}

func defaultFor(forValue string) string {
	if strings.TrimSpace(forValue) == "" {
		return "0m"
	}
	return forValue
}

func normalizeLabelsJSON(labels datatypes.JSON) datatypes.JSON {
	if len(labels) == 0 || string(labels) == "null" {
		return datatypes.JSON([]byte("{}"))
	}
	return labels
}
