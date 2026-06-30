package executor

import (
	"testing"
)

func TestParseRuleGroupFromYAML_alertRules(t *testing.T) {
	content := []byte(`
groups:
  - name: cpu-alerts
    rules:
      - alert: HighCPU
        expr: cpu > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: CPU high
          description: details
          runbook_url: https://example.com
`)
	parsed, err := ParseRuleGroupFromYAML("cpu-alerts.yml", content)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.GroupName != "cpu-alerts" || parsed.GroupType != AlertingRules {
		t.Fatalf("group: name=%q type=%q", parsed.GroupName, parsed.GroupType)
	}
	if len(parsed.Rules) != 1 || parsed.Rules[0].Kind != RuleEntryAlert || parsed.Rules[0].Name != "HighCPU" {
		t.Fatalf("rules: %#v", parsed.Rules)
	}
}

func TestParseRuleGroupsFromYAML_yamlNameOverridesFileName(t *testing.T) {
	content := []byte(`
groups:
  - name: http_uri_rate_aggregation
    rules:
      - record: http:server_requests:rate5m_by_uri
        expr: rate(http[5m])
`)
	groups, err := ParseRuleGroupsFromYAML("http_record.yml", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 || groups[0].GroupName != "http_uri_rate_aggregation" {
		t.Fatalf("groups: %#v", groups)
	}
	if groups[0].GroupType != AlertingRecords {
		t.Fatalf("type=%q", groups[0].GroupType)
	}
}

func TestParseRuleGroupsFromYAML_multipleGroups(t *testing.T) {
	content := []byte(`
groups:
  - name: group_a
    rules:
      - alert: A
        expr: up == 0
        for: 1m
  - name: group_b
    rules:
      - alert: B
        expr: up == 1
        for: 2m
`)
	groups, err := ParseRuleGroupsFromYAML("node_middleware_ruler.yml", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("len=%d", len(groups))
	}
	if groups[0].GroupName != "group_a" || groups[1].GroupName != "group_b" {
		t.Fatalf("names: %q %q", groups[0].GroupName, groups[1].GroupName)
	}
}

func TestParseRuleGroupFromYAML_recordRules(t *testing.T) {
	content := []byte(`
groups:
  - name: agg
    rules:
      - record: instance:cpu:rate5m
        expr: rate(cpu[5m])
`)
	parsed, err := ParseRuleGroupFromYAML("agg.yml", content)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.GroupType != AlertingRecords {
		t.Fatalf("type=%q", parsed.GroupType)
	}
	if parsed.Rules[0].Kind != RuleEntryRecord {
		t.Fatalf("kind=%q", parsed.Rules[0].Kind)
	}
}

func TestParseRuleGroupsFromYAML_mixedRejected(t *testing.T) {
	content := []byte(`
groups:
  - name: mixed
    rules:
      - alert: A
        expr: up
      - record: B:metric
        expr: up
`)
	_, err := ParseRuleGroupsFromYAML("mixed.yml", content)
	if err == nil {
		t.Fatal("expected error for mixed alert and record")
	}
}

func TestParseRuleGroupsFromYAML_singleGroupUsesFileNameWhenNameEmpty(t *testing.T) {
	content := []byte(`
groups:
  - rules:
      - alert: A
        expr: up
`)
	groups, err := ParseRuleGroupsFromYAML("cpu.yml", content)
	if err != nil {
		t.Fatal(err)
	}
	if groups[0].GroupName != "cpu" {
		t.Fatalf("name=%q", groups[0].GroupName)
	}
}
