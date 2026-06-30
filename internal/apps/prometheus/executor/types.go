package executor

import "gorm.io/datatypes"

const (
	AlertingRules   = "ALERTING RULES"
	AlertingRecords = "ALERTING RECORDS"
)

type RuleConfig struct {
	Groups []RuleGroups `yaml:"groups,omitempty" json:"groups,omitempty"`
}

type RuleGroups struct {
	Name  string  `yaml:"name,omitempty" json:"name,omitempty"`
	Rules []Rules `yaml:"rules,omitempty" json:"rules,omitempty"`
}

type Rules struct {
	Expr        string         `yaml:"expr" json:"expr"`
	For         string         `yaml:"for" json:"for"`
	Labels      datatypes.JSON `yaml:"labels" json:"labels"`
	Alert       string         `yaml:"alert" json:"alert"`
	Annotations datatypes.JSON `yaml:"annotations" json:"annotations"`
}

// RecordConfig 结构体
type RecordConfig struct {
	Groups []RecordGroups `yaml:"groups,omitempty" json:"groups,omitempty"`
}
type RecordGroups struct {
	Name  string    `yaml:"name,omitempty" json:"name,omitempty"`
	Rules []Records `yaml:"rules,omitempty" json:"rules,omitempty"`
}

type Records struct {
	Record string         `yaml:"record" json:"record"`
	Expr   string         `yaml:"expr" json:"expr"`
	Labels datatypes.JSON `yaml:"labels" json:"labels,omitempty"`
}

type Rule struct {
	Expr        string
	For         string
	Labels      datatypes.JSON
	Alert       string
	Annotations struct {
		Summary     string
		Description string
	}
}

type Group struct {
	Name  string
	Rules []Rule
}

type Config struct {
	Groups []Group `yaml:"groups"`
}
