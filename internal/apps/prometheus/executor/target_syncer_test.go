package executor

import (
	"encoding/json"
	"testing"

	"prom-lens-backend/internal/apps/prometheus/model"

	"gorm.io/datatypes"
)

func TestBuildTargetGroupJSON_FileSDFormat(t *testing.T) {
	s := &configMapTargetSyncer{}
	enabled := true
	group := model.TargetGroup{
		Name:        "ecs-nodes",
		Description: "非 k8s 节点 node-exporter 采集",
		Labels:      datatypes.JSON(`{"job":"node-exporter"}`),
	}
	targets := []model.Target{
		{
			IPAddress: "10.10.10.10",
			Port:      9100,
			Labels:    datatypes.JSON(`{"hostname":"node01"}`),
			Enabled:   &enabled,
		},
	}

	data, err := s.buildTargetGroupJSON(group, targets)
	if err != nil {
		t.Fatal(err)
	}

	var got []fileSDTarget
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v\nraw=%s", err, string(data))
	}
	if len(got) != 1 {
		t.Fatalf("len=%d want 1", len(got))
	}
	if len(got[0].Targets) != 1 || got[0].Targets[0] != "10.10.10.10:9100" {
		t.Fatalf("targets=%v", got[0].Targets)
	}
	if got[0].Labels["job"] != "node-exporter" || got[0].Labels["hostname"] != "node01" {
		t.Fatalf("labels=%v", got[0].Labels)
	}
}

func TestBuildTargetGroupJSON_SkipsDisabled(t *testing.T) {
	s := &configMapTargetSyncer{}
	disabled := false
	group := model.TargetGroup{Name: "ecs-nodes", Labels: datatypes.JSON(`{}`)}
	targets := []model.Target{{
		IPAddress: "10.10.10.10",
		Port:      9100,
		Labels:    datatypes.JSON(`{}`),
		Enabled:   &disabled,
	}}

	data, err := s.buildTargetGroupJSON(group, targets)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "[]" {
		t.Fatalf("got %s want []", string(data))
	}
}
