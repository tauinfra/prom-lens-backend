package prometheus

import (
	"encoding/json"
	"testing"

	"gorm.io/datatypes"
)

func TestValidateExtraAnnotations_reservedKey(t *testing.T) {
	extra := datatypes.JSON([]byte(`{"summary":"x"}`))
	if err := ValidateExtraAnnotations(extra); err == nil {
		t.Fatal("expected error for reserved key summary")
	}
}

func TestMergeRuleAnnotations(t *testing.T) {
	extra := datatypes.JSON([]byte(`{"runbook_url":"https://example.com","owner":"sre"}`))
	merged, err := MergeRuleAnnotations("cpu high", "instance {{ $labels.instance }}", extra)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	if err := json.Unmarshal(merged, &m); err != nil {
		t.Fatal(err)
	}
	if m["summary"] != "cpu high" || m["description"] != "instance {{ $labels.instance }}" {
		t.Fatalf("unexpected fixed keys: %#v", m)
	}
	if m["runbook_url"] != "https://example.com" || m["owner"] != "sre" {
		t.Fatalf("unexpected extra keys: %#v", m)
	}
}

func TestSplitRuleAnnotations(t *testing.T) {
	merged := datatypes.JSON([]byte(`{"summary":"s","description":"d","runbook_url":"https://x"}`))
	summary, desc, extra, err := SplitRuleAnnotations(merged)
	if err != nil {
		t.Fatal(err)
	}
	if summary != "s" || desc != "d" {
		t.Fatalf("got summary=%q description=%q", summary, desc)
	}
	var m map[string]string
	_ = json.Unmarshal(extra, &m)
	if m["runbook_url"] != "https://x" || len(m) != 1 {
		t.Fatalf("extra: %#v", m)
	}
}

func TestMergeRuleAnnotations_ignoresReservedInExtra(t *testing.T) {
	extra := datatypes.JSON([]byte(`{"summary":"override","runbook_url":"https://example.com"}`))
	merged, err := MergeRuleAnnotations("from column", "desc", extra)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]string
	_ = json.Unmarshal(merged, &m)
	if m["summary"] != "from column" {
		t.Fatalf("summary should come from columns, got %q", m["summary"])
	}
	if m["runbook_url"] != "https://example.com" {
		t.Fatalf("extra key missing: %#v", m)
	}
}
