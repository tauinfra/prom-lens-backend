package helper

import (
	"fmt"
	"strings"
)

// EventHelper 节点相关工具函数
type EventHelper struct{}

func NewEventHelper() *EventHelper {
	return &EventHelper{}
}

// BuildEventFieldSelector 生成 Event 事件过滤标签
func (h *EventHelper) BuildEventFieldSelector(kind, name string) string {
	var fieldSelectors []string
	if kind != "" {
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("involvedObject.kind=%s", kind))
	}
	if name != "" {
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("involvedObject.name=%s", name))
	}
	return strings.Join(fieldSelectors, ",")
}
