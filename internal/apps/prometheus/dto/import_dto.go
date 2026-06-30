package dto

// ImportRulesResult 全量 ConfigMap 导入结果。
type ImportRulesResult struct {
	FilesTotal       int              `json:"filesTotal"`
	GroupsCreated    int              `json:"groupsCreated"`
	GroupsSkipped    int              `json:"groupsSkipped"`
	RulesCreated     int              `json:"rulesCreated"`
	RulesSkipped     int              `json:"rulesSkipped"`
	RecordsCreated   int              `json:"recordsCreated"`
	RecordsSkipped   int              `json:"recordsSkipped"`
	Errors           []ImportErrorDTO `json:"errors"`
}

type ImportErrorDTO struct {
	ConfigMapKey string `json:"configMapKey,omitempty"`
	GroupName    string `json:"groupName,omitempty"`
	ItemName     string `json:"itemName,omitempty"`
	Message      string `json:"message"`
}
